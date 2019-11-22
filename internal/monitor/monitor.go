package monitor

import (
	"fmt"
	"sync"
	"time"

	"github.com/hpifu/go-monitor/internal/collector"
	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	influxdb "github.com/influxdata/influxdb1-client/v2"
	"github.com/sirupsen/logrus"
)

func NewMonitor(addr string, database string, batch int, retry int) (*Monitor, error) {
	c, err := influxdb.NewHTTPClient(influxdb.HTTPConfig{
		Addr: addr,
	})
	if err != nil {
		return nil, err
	}

	if _, _, err := c.Ping(200 * time.Millisecond); err != nil {
		return nil, err
	}

	return &Monitor{
		client:   c,
		database: database,

		metricQueue: make(chan *metricItem, 1000),
		stop:        false,

		batch:     batch,
		retry:     retry,
		infoLog:   logrus.New(),
		warnLog:   logrus.New(),
		accessLog: logrus.New(),
	}, nil
}

type metricItem struct {
	table  string
	metric *collector.Metric
}

type Monitor struct {
	client   influxdb.Client
	database string

	metricQueue chan *metricItem

	collectorWG sync.WaitGroup
	monitorWG   sync.WaitGroup
	stop        bool

	batch int
	retry int

	infoLog   *logrus.Logger
	warnLog   *logrus.Logger
	accessLog *logrus.Logger
}

func (m *Monitor) SetLogger(infoLog, warnLog, accessLog *logrus.Logger) {
	m.infoLog = infoLog
	m.warnLog = warnLog
	m.accessLog = accessLog
}

func (m *Monitor) AddCollector(c collector.Collector, table string, interval time.Duration) {
	go func() {
		m.infoLog.Infof("add collector, table: [%v], interval: [%v]", table, interval)
		m.collectorWG.Add(1)
		for range time.Tick(interval) {
			for _, metric := range c.Collect() {
				m.metricQueue <- &metricItem{
					table:  table,
					metric: metric,
				}
			}
			if m.stop {
				break
			}
		}
		m.collectorWG.Done()
		m.infoLog.Infof("collector done")
	}()
}

func (m *Monitor) Monitor() error {
	var bps influxdb.BatchPoints
	var err error
	count := 0

	bps, err = influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
		Database:  m.database,
		Precision: "s",
	})
	if err != nil {
		return err
	}

	go func() {
		m.infoLog.Infof("monitor start")
		m.monitorWG.Add(1)
		for item := range m.metricQueue {
			fmt.Println(item.table, item.metric)
			point, err := influxdb.NewPoint(
				item.table, item.metric.Keys, item.metric.Vals, item.metric.Timestamp,
			)
			if err != nil {
				m.warnLog.Warnf("new point failed. err: [%v]", err)
				continue
			}

			m.accessLog.WithFields(logrus.Fields{
				"table":  item.table,
				"metric": item.metric,
			}).Info()
			bps.AddPoint(point)
			count++

			if count == m.batch {
				retry := 0
				for err = m.client.Write(bps); err != nil; {
					m.warnLog.Warnf("influxdb write failed. err: [%v], retry: [%v]", err, retry)
					retry++
					if retry == m.retry {
						break
					}
				}

				bps, _ = influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
					Database:  m.database,
					Precision: "s",
				})
				count = 0
			}
		}

		if count != 0 {
			retry := 0
			for err = m.client.Write(bps); err != nil; {
				m.warnLog.Warnf("influxdb write failed. err: [%v], retry: [%v]", err, retry)
				retry++
				if retry == m.retry {
					break
				}
			}
		}

		m.monitorWG.Done()
		m.infoLog.Infof("monitor done")
	}()

	return nil
}

func (m *Monitor) Stop() {
	m.stop = true
	m.collectorWG.Wait()
	close(m.metricQueue)
	m.monitorWG.Wait()

	_ = m.client.Close()
}
