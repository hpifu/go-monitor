package collector

import "syscall"

type DiskCollector struct{}

func NewDiskCollector() (*DiskCollector, error) {
	return &DiskCollector{}, nil
}

func (c *DiskCollector) Collect() map[string]float64 {
	stats := &syscall.Statfs_t{}
	_ = syscall.Statfs("/home", stats)

	return map[string]float64{
		"total": float64(stats.Blocks*uint64(stats.Bsize)) / GBytes,
		"free":  float64(stats.Bfree*uint64(stats.Bsize)) / GBytes,
		"used":  float64((stats.Blocks-stats.Bfree)*uint64(stats.Bsize)) / GBytes,
	}
}
