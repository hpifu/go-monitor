# 系统监控

监控宿主机器的一些基础指标，并写入到 influxdb 用于之后的可视化以及报警服务

监控的指标包括：cpu利用率，cpu负载，内存使用，网络负载，iops，磁盘等

## 总体设计

![总体架构](go-monitor.png)

主要分为两大模块，`reporter`，`collector` 和 `scheduler`

- `collector`: 负责具体的某些监控指标的采集
- `reporter`: 负责将采集到的指标写入到数据库
- `scheduler`: 负责整个流程的调度，数据采集和数据写入的协同

`collector` 将采集到的数据写到 `channel` 里，`monitor` 从 `channel` 中读取采集到的数据批量写入到 `influxdb` 中，整个过程的系统由 `scheduler` 来调度

## 设计思路

### collector 设计

``` go
type Metric struct {
	Keys      map[string]string      `json:"keys"`
	Vals      map[string]interface{} `json:"vals"`
	Timestamp time.Time              `json:"timestamp"`
}
```

metric 是指标数据的抽象，包含三个字段

- `Keys`: 指标的维度，对应 influxdb 中的 tag，主要用于查询时对数据进行分类，目前只有 host 在该字段中，后面可以按需添加新的字段
- `Vals`: 指标的值，对应 influxdb 中的 field，用于记录指标具体的值，一条数据中可以有多个指标，比如 cpu 利用率就有 system/user/idle 三个值，`Vals` 被设计成 `map[string]interface{}` 主要是为了和 influxdb 提供的接口对其，目前只有 float64 类型
- `Timestamp`: 指标采集的时间，由于数据是批量发送，刚生成的数据可能会等待下一条数据一起打包发送，这个时间间隔可能较长，因此在每条数据里面加上这条数据产生的时间戳，一起发送给 influxdb

``` go
type Collector interface {
	Collect() []*Metric
}
```

所有的数据采集过程抽象成一个 `collector` 接口，在整个工作流中，这个接口会被周期性调用，每次调用返回一条或多条 metric 数据(目前仅有一条，但可见的拓展需求，比如多网卡或多磁盘的监控就可能返回多条数据)

``` go
func NewCollector(name string, params []interface{}) (Collector, error)
```

再提供一个工厂方法，通过类名和参数来构造 collector

### reporter 设计

``` go
type MetricItem struct {
	Table  string
	Metric *collector.Metric
}

type Reporter interface {
	Report([]*MetricItem) error
}
```

具体某条数据的写入到哪个表中，由 `scheduler` 通过 `MetricItem` 告知 `reporter`

`reporter` 提供一个 `Report` 供 scheduler 调用，将 mertic 数据打包写入到数据库

### scheduler 设计

``` go
func (s *Scheduler) SetReporter(r reporter.Reporter)
func (s *Scheduler) AddCollector(c collector.Collector, table string, interval time.Duration)
func (s *Scheduler) Scheduler()
func (s *Scheduler) Stop()
```

主要提供三个接口

- `SetReporter`: 设置 `reporter`，目前只有 influxdb 一个 `reporter`，目前一个 `scheduler` 中只有一个 `reporter`
- `AddCollector`: 新增 `collector`，一个 `scheduler` 可以有多个 `collecotr`
- `Scheduler`: 开始调度，这里会创建数据采集和数据写入协程，通过 channel 实现协程之间的通信
- `Stop`: 调度结束，停止所有数据采集协程，发送队列中所有剩余的数据，然后退出

### 主要工作流程

1. 通过配置分别构造 `collector`，`reporter`，`scheduler` 对象
2. 调用 `scheduler.AddCollector`，`scheduler.SetReporter` 将 `collector` 和 `reporter` 对象关联到 `scheduler` 中
3. 调用 `scheduler.Scheduler` 开始调度，这里将创建周期执行的 collector 协程，以及负责数据写入的 reporter 协程
4. 等待退出信号，退出时，先停止当前的 collector 协程，再等待 reporter 协程退出

``` go
c := NewCollector()
r := NewReporter()
s := NewScheduler()

s.AddCollector(c)
s.SetReporter(r)
s.Scheduler()

s.Stop()
```

### reporter 和 collector 的同步

这是一个典型的生成者消费者问题，多个生成者 `collector` 往一个 `channel` 中写入，一个消费者 `reporter` 从 `channel` 中消费数据

需要注意的是，执行退出时，需要将队列中数据消费完，正确的执行顺序应该为：

1. `collector` 停止写入
2. 等待所有的 `collector` 退出
3. 关闭 `channel`
4. 等待 `reporter` 退出

``` go
s.stop = true
s.collectorWG.Wait()
close(s.metricQueue)
s.reporterWG.Wait()
```

这里使用两个 `sync.WaitGroup` 来同步，分别用于等待 `collector` 和 `reporter`

### 数据采集

目前主要有 cpu利用率，cpu负载，内存使用，网络负载，iops，磁盘这些数据的采集，主要使用 `github.com/mackerelio/go-osstat` 相关接口，这个库对各个操作系统的监控作了封装，并提供了统一的即可，linux 下实现基本都是解析 `/proc` 目录下系统文件的数据

**cpu 利用率**

`github.com/mackerelio/go-osstat/cpu` 下的 `Get` 接口返回当前总的 cpu 时间(user/system/idle)，需要在每次调用减去上一次调用的值，可以得到这段时间之内 cpu 时间，这段时间内 user/system/idle 与 total 的比值就是 cpu 利用率

结果应该与 `top` 命令观察结果一致

**cpu 负载**

`github.com/mackerelio/go-osstat/loadavg` 下的 `Get` 接口返回 1分钟，5分钟，15分钟内的平均负载，因此直接返回这个接口即可

结果应该与 `uptime` 命令观察结果一致

**内存使用**

`github.com/mackerelio/go-osstat/memory` 下的 `Get` 接口返回当前内存的使用情况，直接返回这个接口即可

结果应该与 `free` 命令观察结果一致

**网络负载**

`github.com/mackerelio/go-osstat/network` 下的 `Get` 接口返回各个网卡下的网络流量（包括一些虚拟网卡），这里我们只关注外网的流量，这个设备名一般是 eth0，centos7 为了支持多个网卡设备名的唯一性，改成了以 `en` 开头的网卡，因此我们需要从返回的数据中找到以 `en` 开头的那个那个网卡，返回对应的数据

结果应该和 `netstat -i` 命令观察结果一致

**iops**

`github.com/mackerelio/go-osstat/disk` 下的 `Get` 接口返回各个磁盘设备总的 io 次数，这里我们服务的磁盘，一般是 `sda`，因此需要先找到 `sda`，再减去上次的 io 次数，除以时间得到 iops

**磁盘大小**

磁盘大小监控直接使用 golang 系统的 api 即可

``` golang
stats := &syscall.Statfs_t{}
_ = syscall.Statfs("/home", stats)
```
