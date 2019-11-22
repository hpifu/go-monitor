package collector

import (
	"github.com/mackerelio/go-osstat/memory"
	"time"
)

// linux command: free
type MemoryCollector struct {
	keys map[string]string
}

func NewMemoryCollector() (*MemoryCollector, error) {
	return &MemoryCollector{
		keys: map[string]string{
			"host": Hostname(),
		},
	}, nil
}

func (c *MemoryCollector) Collect() []*Metric {
	value, _ := memory.Get()

	return []*Metric{
		{
			Keys: c.keys,
			Vals: map[string]interface{}{
				"total":  float64(value.Total) / GBytes,
				"free":   float64(value.Free) / GBytes,
				"used":   float64(value.Used) / GBytes,
				"cached": float64(value.Cached) / GBytes,
			},
			Timestamp: time.Now(),
		},
	}
}
