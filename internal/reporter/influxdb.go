package reporter

import (
	"time"

	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	influxdb "github.com/influxdata/influxdb1-client/v2"
)

func NewInfluxdbReporter(addr string, database string, retry int) (*InfluxdbReporter, error) {
	c, err := influxdb.NewHTTPClient(influxdb.HTTPConfig{
		Addr: addr,
	})
	if err != nil {
		return nil, err
	}

	if _, _, err := c.Ping(200 * time.Millisecond); err != nil {
		return nil, err
	}

	return &InfluxdbReporter{
		client:   c,
		database: database,
		retry:    retry,
	}, nil
}

type InfluxdbReporter struct {
	client   influxdb.Client
	database string
	retry    int
}

func (m *InfluxdbReporter) Report(items []*MetricItem) error {
	bps, err := influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
		Database:  m.database,
		Precision: "s",
	})
	if err != nil {
		return err
	}

	for _, item := range items {
		point, err := influxdb.NewPoint(
			item.Table, item.Metric.Keys, item.Metric.Vals, item.Metric.Timestamp,
		)
		if err != nil {
			return err
		}
		bps.AddPoint(point)

		retry := 0
		for err = m.client.Write(bps); err != nil; {
			retry++
			if retry == m.retry {
				return err
			}
		}

		bps, _ = influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
			Database:  m.database,
			Precision: "s",
		})
	}

	return nil
}
