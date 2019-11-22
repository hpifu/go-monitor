package collector

import (
	"os"
	"time"
)

const KBytes = float64(1024)
const MBytes = float64(KBytes * 1024)
const GBytes = float64(MBytes * 1024)
const Mbytes = float64(MBytes / 8)

type Metric struct {
	Keys      map[string]string
	Vals      map[string]interface{}
	Timestamp time.Time
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
