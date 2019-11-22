package collector

import (
	"fmt"
	"os"
	"time"
)

const KBytes = float64(1024)
const MBytes = float64(KBytes * 1024)
const GBytes = float64(MBytes * 1024)
const Mbytes = float64(MBytes / 8)

type Metric struct {
	Keys      map[string]string      `json:"keys"`
	Vals      map[string]interface{} `json:"vals"`
	Timestamp time.Time              `json:"timestamp"`
}

type Collector interface {
	Collect() []*Metric
}

func Hostname() string {
	host, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return host
}

func NewCollector(name string, params []interface{}) (Collector, error) {
	if name == "CPUCollector" {
		return NewCPUCollector()
	}
	if name == "MemoryCollector" {
		return NewMemoryCollector()
	}
	if name == "NetworkCollector" {
		if len(params) == 0 {
			return nil, fmt.Errorf("params required")
		}
		if _, ok := params[0].(string); !ok {
			return nil, fmt.Errorf("params should be string, [%v]", params[0])
		}
		return NewNetworkCollector(params[0].(string))
	}
	if name == "LoadAvgCollector" {
		return NewLoadAvgCollector()
	}
	if name == "DiskCollector" {
		return NewDiskCollector()
	}
	if name == "IOPSCollector" {
		if len(params) == 0 {
			return nil, fmt.Errorf("params required")
		}
		if _, ok := params[0].(string); !ok {
			return nil, fmt.Errorf("params should be string, [%v]", params[0])
		}
		//return NewIOPSCollector(params[0].(string))
		return NewDiskCollector()
	}

	return nil, fmt.Errorf("no such collector")
}
