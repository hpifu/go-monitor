package collector

import (
	"github.com/mackerelio/go-osstat/loadavg"
	"time"
)

// linux command: uptime
type LoadAvgCollector struct {
	keys map[string]string
}

func NewLoadAvgCollector() (*LoadAvgCollector, error) {
	return &LoadAvgCollector{
		keys: map[string]string{
			"host": Hostname(),
		},
	}, nil
}

func (c *LoadAvgCollector) Collect() []*Metric {
	value, _ := loadavg.Get()

	return []*Metric{
		{
			Keys: c.keys,
			Vals: map[string]interface{}{
				"load1m":  value.Loadavg1,
				"load5m":  value.Loadavg5,
				"load15m": value.Loadavg15,
			},
			Timestamp: time.Now(),
		},
	}
}
