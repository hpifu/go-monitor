package monitor

import (
	"fmt"
	"github.com/hpifu/go-monitor/internal/collector"
	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	influxdb "github.com/influxdata/influxdb1-client/v2"
	"sync"
	"time"
)

type metricItem struct {
	table  string
	metric *collector.Metric
}

type Monitor struct {
	client influxdb.Client
	dbname string

	metricQueue chan *metricItem

	collectorWG sync.WaitGroup
	monitorWG   sync.WaitGroup
	stop        bool

	batch int
	retry int
}

func NewMonitor(addr string, dbname string) (*Monitor, error) {
	c, err := influxdb.NewHTTPClient(influxdb.HTTPConfig{
		Addr: addr,
	})

	if err != nil {
		return nil, err
	}

	return &Monitor{
		client: c,
		dbname: dbname,

		metricQueue: make(chan *metricItem, 1000),
		stop:        false,

		batch: 5,
		retry: 3,
	}, nil
}

func (m *Monitor) AddCollector(c collector.Collector, table string, interval time.Duration) {
	go func() {
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
	}()
}

func (m *Monitor) Monitor() error {
	var bps influxdb.BatchPoints
	var err error
	count := 0

	bps, err = influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
		Database:  m.dbname,
		Precision: "s",
	})
	if err != nil {
		return err
	}

	go func() {
		m.monitorWG.Add(1)
		for item := range m.metricQueue {
			fmt.Println(item.table, item.metric)
			point, err := influxdb.NewPoint(
				item.table, item.metric.Keys, item.metric.Vals, item.metric.Timestamp,
			)
			if err != nil {
				// do some log
				continue
			}

			bps.AddPoint(point)
			count++

			if count == m.batch {
				retry := 0
				for err = m.client.Write(bps); err != nil; {
					// do some log
					retry++
					if retry == m.retry {
						break
					}
				}

				bps, _ = influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
					Database:  m.dbname,
					Precision: "s",
				})
				count = 0
			}
		}

		if count != 0 {
			retry := 0
			for err = m.client.Write(bps); err != nil; {
				// do some log
				retry++
				if retry == m.retry {
					break
				}
			}
		}

		m.monitorWG.Done()
	}()

	return nil
}

func (m *Monitor) Save(metric *collector.Metric) error {
	bps, err := influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
		Database:  m.dbname,
		Precision: "s",
	})
	if err != nil {
		return err
	}

	point, err := influxdb.NewPoint("mydb", metric.Keys, metric.Vals, metric.Timestamp)

	if err != nil {
		return err
	}

	bps.AddPoint(point)

	return m.client.Write(bps)
}

func (m *Monitor) Stop() {
	m.stop = true
	m.collectorWG.Wait()
	close(m.metricQueue)
	m.monitorWG.Wait()

	_ = m.client.Close()
}
