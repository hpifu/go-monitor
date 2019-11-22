package collector

import "github.com/mackerelio/go-osstat/loadavg"

// linux command: uptime
type LoadAvgCollector struct {}

func NewLoadAvgCollector() (*LoadAvgCollector, error) {
	return &LoadAvgCollector{}, nil
}

func (c *LoadAvgCollector) Collect() map[string]float64 {
	value, _ := loadavg.Get()

	return map[string]float64 {
		"load1m": value.Loadavg1,
		"load5m": value.Loadavg5,
		"load15m": value.Loadavg15,
	}
}
