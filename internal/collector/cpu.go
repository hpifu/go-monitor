package collector

import "github.com/mackerelio/go-osstat/cpu"

type CPUCollector struct {
	value *cpu.Stats
}

func NewCPUCollector() (*CPUCollector, error) {
	value, err := cpu.Get()
	if err != nil {
		return nil, err
	}

	return &CPUCollector{
		value: value,
	}, nil
}

func (c *CPUCollector) Collect() map[string]float64 {
	value, _ := cpu.Get()
	total := float64(value.Total - c.value.Total)
	user := float64(value.User - c.value.User)
	system := float64(value.System - c.value.System)
	idel := float64(value.Idle - c.value.Idle)

	c.value = value

	return map[string]float64{
		"user": user / total * 100,
		"collector": system / total * 100,
		"idel": idel / total * 100,
	}
}
