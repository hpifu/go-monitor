package collector

import (
	"github.com/mackerelio/go-osstat/cpu"
	"math"
	"time"
)

// linux command: top
type CPUCollector struct {
	value *cpu.Stats
	keys  map[string]string
}

func NewCPUCollector() (*CPUCollector, error) {
	value, err := cpu.Get()
	if err != nil {
		return nil, err
	}

	return &CPUCollector{
		keys: map[string]string{
			"host": Hostname(),
		},
		value: value,
	}, nil
}

func (c *CPUCollector) Collect() []*Metric {
	value, _ := cpu.Get()
	total := float64(value.Total - c.value.Total)
	user := float64(value.User - c.value.User)
	system := float64(value.System - c.value.System)
	idle := float64(value.Idle - c.value.Idle)

	c.value = value

	return []*Metric{
		{
			Keys: c.keys,
			Vals: map[string]interface{}{
				"user":   math.Round(user/total*10000) / 100,
				"system": math.Round(system/total*10000) / 100,
				"idle":   math.Round(idle/total*10000) / 100,
			},
			Timestamp: time.Now(),
		},
	}
}
