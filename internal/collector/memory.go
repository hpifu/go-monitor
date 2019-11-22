package collector

import "github.com/mackerelio/go-osstat/memory"

type MemoryCollector struct {
	value *memory.Stats
}

func NewMemoryCollector() (*MemoryCollector, error) {
	value, err := memory.Get()
	if err != nil {
		return nil, err
	}

	return &MemoryCollector{
		value: value,
	}, nil
}

func (c *MemoryCollector) Collect() map[string]float64 {
	value, _ := memory.Get()

	return map[string]float64 {
		"total": float64(value.Total),
		"free": float64(value.Free),
		"used": float64(value.Used),
		"cached": float64(value.Cached),
	}
}
