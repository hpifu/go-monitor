package collector

import (
	"github.com/mackerelio/go-osstat/loadavg"
	"math"
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
				"load1m":  math.Round(value.Loadavg1*100) / 100,
				"load5m":  math.Round(value.Loadavg5*100) / 100,
				"load15m": math.Round(value.Loadavg15*100) / 100,
			},
			Timestamp: time.Now(),
		},
	}
}
