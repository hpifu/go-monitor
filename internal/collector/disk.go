package collector

import (
	"math"
	"syscall"
	"time"
)

// linux command: df -h
type DiskCollector struct {
	keys map[string]string
}

func NewDiskCollector() (*DiskCollector, error) {
	return &DiskCollector{
		keys: map[string]string{
			"host": Hostname(),
		},
	}, nil
}

func (c *DiskCollector) Collect() []*Metric {
	stats := &syscall.Statfs_t{}
	_ = syscall.Statfs("/home", stats)

	return []*Metric{
		{
			Keys: c.keys,
			Vals: map[string]interface{}{
				"total": math.Round(float64(stats.Blocks*uint64(stats.Bsize))/GBytes*100) / 100,
				"free":  math.Round(float64(stats.Bfree*uint64(stats.Bsize))/GBytes*100) / 100,
				"used":  math.Round(float64((stats.Blocks-stats.Bfree)*uint64(stats.Bsize))/GBytes*100) / 100,
			},
			Timestamp: time.Now(),
		},
	}
}
