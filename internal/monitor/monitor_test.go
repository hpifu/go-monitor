package monitor

import (
	"testing"
	"time"

	"github.com/hpifu/go-monitor/internal/collector"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMonitor(t *testing.T) {
	Convey("test monitor", t, func() {
		m, err := NewMonitor("http://localhost:8086", "monitor", 5, 3)
		So(err, ShouldBeNil)

		cpuCollector, _ := collector.NewCPUCollector()
		memCollector, _ := collector.NewMemoryCollector()
		m.AddCollector(cpuCollector, "cpu", time.Second)
		m.AddCollector(memCollector, "mem", time.Second)
		So(m.Monitor(), ShouldBeNil)
		time.Sleep(13 * time.Second)
		m.Stop()
	})
}
