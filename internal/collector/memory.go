package collector

import (
	"github.com/mackerelio/go-osstat/memory"
	"math"
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
				"total":  math.Round(float64(value.Total)/GBytes*100) / 100,
				"free":   math.Round(float64(value.Free)/GBytes*100) / 100,
				"used":   math.Round(float64(value.Used)/GBytes*100) / 100,
				"cached": math.Round(float64(value.Cached)/GBytes*100) / 100,
			},
			Timestamp: time.Now(),
		},
	}
}
