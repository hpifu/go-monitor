package reporter

import "github.com/hpifu/go-monitor/internal/collector"

type MetricItem struct {
	Table  string
	Metric *collector.Metric
}

type Reporter interface {
	Report([]*MetricItem) error
}
