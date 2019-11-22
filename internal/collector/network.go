package collector

import (
	"github.com/mackerelio/go-osstat/memory"
	"github.com/mackerelio/go-osstat/network"
)

type NetworkCollector struct {
	value *memory.Stats
}

func NewNetworkCollector() (*NetworkCollector, error) {
	value, err := memory.Get()
	if err != nil {
		return nil, err
	}

	return &NetworkCollector{
		value: value,
	}, nil
}

func (c *NetworkCollector) Collect() map[string]float64 {
	value, _:= network.Get()

	return map[string]float64 {
		//"total": float64(value.),
	}
}
