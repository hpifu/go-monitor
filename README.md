# 系统监控

监控宿主机器的一些基础指标，并写入到 influxdb 用于之后的可视化以及报警服务

监控的指标包括：cpu利用率，cpu负载，内存使用率，网络负载，iops，磁盘等

## 总体设计

![总体架构](go-monitor.png)

主要分为两大模块，`monitor` 和 `collector`

- `collector`: 负责具体的某些监控指标的采集
- `monitor`: 负责将采集到的指标写入到数据库

collector 将采集到的数据写到 channel 里，monitor 从 channel 中读取采集到的数据批量写入到 influxdb 中

## 设计思路

### metric 数据结构

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

### collector 接口

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

### monitor 设计

``` go
func (m *Monitor) AddCollector(c collector.Collector, table string, interval time.Duration)
func (m *Monitor) Monitor() error
```

### 主要工作流程



### monitor 和 collector 的同步
