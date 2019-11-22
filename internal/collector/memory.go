package collector

import "github.com/mackerelio/go-osstat/memory"

// linux command: free
type MemoryCollector struct{}

func NewMemoryCollector() (*MemoryCollector, error) {
	return &MemoryCollector{}, nil
}

func (c *MemoryCollector) Collect() map[string]float64 {
	value, _ := memory.Get()

	return map[string]float64{
		"total":  float64(value.Total) / GBytes,
		"free":   float64(value.Free) / GBytes,
		"used":   float64(value.Used) / GBytes,
		"cached": float64(value.Cached) / GBytes,
	}
}
