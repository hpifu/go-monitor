package collector

import (
	"fmt"
	"strings"
	"time"

	"github.com/mackerelio/go-osstat/disk"
)

// linux command: iostat
type IOPSCollector struct {
	name  string
	value *disk.Stats
	ts    time.Time
}

func NewIOPSCollector(name string) (*IOPSCollector, error) {
	value := collectIOPS(name)
	ts := time.Now()

	if value == nil {
		return nil, fmt.Errorf("network not found, name: [%v]", name)
	}

	return &IOPSCollector{
		name:  name,
		value: value,
		ts:    ts,
	}, nil
}

func collectIOPS(name string) *disk.Stats {
	vals, _ := disk.Get()
	for _, val := range vals {
		if strings.HasPrefix(val.Name, name) {
			return &val
		}
	}

	return nil
}

func (c *IOPSCollector) Collect() map[string]float64 {
	value := collectIOPS(c.name)
	ts := time.Now()

	res := map[string]float64{
		"RMbps": float64(value.ReadsCompleted-c.value.ReadsCompleted) / float64(ts.Sub(c.ts)/time.Second),
		"WMbps": float64(value.WritesCompleted-c.value.WritesCompleted) / float64(ts.Sub(c.ts)/time.Second),
	}

	c.value = value
	c.ts = ts

	return res
}
