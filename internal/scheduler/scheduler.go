package scheduler

import (
	"github.com/hpifu/go-monitor/internal/reporter"
	"sync"
	"time"

	"github.com/hpifu/go-monitor/internal/collector"
	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	"github.com/sirupsen/logrus"
)

func NewScheduler(batch int) *Scheduler {
	return &Scheduler{
		metricQueue: make(chan *reporter.MetricItem, 1000),
		stop:        false,
		batch:       batch,
		infoLog:     logrus.New(),
		warnLog:     logrus.New(),
		accessLog:   logrus.New(),
	}
}

type collectorInfo struct {
	collector collector.Collector
	table     string
	interval  time.Duration
}

type Scheduler struct {
	collectors  []*collectorInfo
	reporter    reporter.Reporter
	metricQueue chan *reporter.MetricItem
	batch       int

	collectorWG sync.WaitGroup
	reporterWG  sync.WaitGroup
	stop        bool

	infoLog   *logrus.Logger
	warnLog   *logrus.Logger
	accessLog *logrus.Logger
}

func (s *Scheduler) SetLogger(infoLog, warnLog, accessLog *logrus.Logger) {
	s.infoLog = infoLog
	s.warnLog = warnLog
	s.accessLog = accessLog
}

func (s *Scheduler) SetReporter(r reporter.Reporter) {
	s.reporter = r
}

func (s *Scheduler) AddCollector(c collector.Collector, table string, interval time.Duration) {
	s.collectors = append(s.collectors, &collectorInfo{
		collector: c,
		table:     table,
		interval:  interval,
	})
}

func (s *Scheduler) Scheduler() error {
	for _, info := range s.collectors {
		go func(info *collectorInfo) {
			s.infoLog.Infof("add collector, table: [%v], interval: [%v]", info.table, info.interval)
			s.collectorWG.Add(1)
			for range time.Tick(info.interval) {
				for _, metric := range info.collector.Collect() {
					s.metricQueue <- &reporter.MetricItem{
						Table:  info.table,
						Metric: metric,
					}
				}
				if s.stop {
					break
				}
			}
			s.collectorWG.Done()
			s.infoLog.Infof("collector done")
		}(info)
	}

	go func() {
		s.infoLog.Infof("Scheduler start")
		s.reporterWG.Add(1)
		var metrics []*reporter.MetricItem
		for item := range s.metricQueue {
			s.accessLog.WithFields(logrus.Fields{
				"table":  item.Table,
				"metric": item.Metric,
			}).Info()
			metrics = append(metrics, item)

			if len(metrics) == s.batch {
				if err := s.reporter.Report(metrics); err != nil {
					s.warnLog.Warnf("influxdb write failed. err: [%v]", err)
				}
				metrics = metrics[:0]
			}
		}

		if len(metrics) != 0 {
			if err := s.reporter.Report(metrics); err != nil {
				s.warnLog.Warnf("influxdb write failed. err: [%v]", err)
			}
		}

		s.reporterWG.Done()
		s.infoLog.Infof("Scheduler done")
	}()

	return nil
}

func (s *Scheduler) Stop() {
	s.stop = true
	s.collectorWG.Wait()
	close(s.metricQueue)
	s.reporterWG.Wait()
}
