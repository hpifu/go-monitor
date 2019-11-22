package collector

import (
	"fmt"
	"github.com/mackerelio/go-osstat/network"
	"math"
	"strings"
	"time"
)

// linux command: netstat -i
type NetworkCollector struct {
	name  string
	value *network.Stats
	ts    time.Time

	keys map[string]string
}

func NewNetworkCollector(name string) (*NetworkCollector, error) {
	value := collectNetwork(name)
	ts := time.Now()

	if value == nil {
		return nil, fmt.Errorf("network not found, name: [%v]", name)
	}

	return &NetworkCollector{
		name:  name,
		value: value,
		ts:    ts,
		keys: map[string]string{
			"host": Hostname(),
		},
	}, nil
}

func collectNetwork(name string) *network.Stats {
	vals, _ := network.Get()
	for _, val := range vals {
		if strings.HasPrefix(val.Name, name) {
			return &val
		}
	}

	return nil
}

func (c *NetworkCollector) Collect() []*Metric {
	value := collectNetwork(c.name)
	ts := time.Now()

	res := map[string]interface{}{
		"imbps": math.Round(float64(value.RxBytes-c.value.RxBytes)/float64(ts.Sub(c.ts)/time.Second)/Mbytes*100) / 100,
		"ombps": math.Round(float64(value.TxBytes-c.value.TxBytes)/float64(ts.Sub(c.ts)/time.Second)/Mbytes*100) / 100,
	}

	c.value = value
	c.ts = ts

	return []*Metric{
		{
			Keys:      c.keys,
			Vals:      res,
			Timestamp: time.Now(),
		},
	}
}
