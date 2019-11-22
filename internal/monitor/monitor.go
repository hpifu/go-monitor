package monitor

import (
	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	influxdb "github.com/influxdata/influxdb1-client/v2"
)

type Mertic struct {
	Table string
	Keys  map[string]string
	Value float64
}

type Monitor struct {
	client influxdb.Client
}

type Collector interface {
	Collect() map[string]float64
}

func NewMonitor(addr string) (*Monitor, error) {
	c, err := influxdb.NewHTTPClient(influxdb.HTTPConfig{
		Addr: addr,
	})

	if err != nil {
		return nil, err
	}

	return &Monitor{
		client: c,
	}, nil
}

func (m *Monitor) Save(mertic *Mertic) error {
	bps, err := influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
		Database:  "mydb",
		Precision: "s",
	})
	if err != nil {
		return err
	}

	point, err := influxdb.NewPoint(mertic.Table, mertic.Keys, map[string]interface{}{
		"value": mertic.Value,
	})

	if err != nil {
		return err
	}

	bps.AddPoint(point)

	return m.client.Write(bps)
}

func (m *Monitor) Close() {
	_ = m.client.Close()
}
