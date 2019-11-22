package monitor

import (
	"github.com/hpifu/go-monitor/internal/collector"
	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	influxdb "github.com/influxdata/influxdb1-client/v2"
)

type Monitor struct {
	client influxdb.Client
	dbname string
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
	}, nil
}

func (m *Monitor) AddCollector(c collector.Collector) {

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

func (m *Monitor) Close() {
	_ = m.client.Close()
}
