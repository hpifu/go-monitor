package collector

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCPUCollector(t *testing.T) {
	Convey("test cpu collector", t, func() {
		c, err := NewCPUCollector()
		So(err, ShouldBeNil)
		for i := 0; i < 10; i++ {
			time.Sleep(time.Second)
			for _, m := range c.Collect() {
				fmt.Println(m)
			}
		}
	})
}

func TestMemoryCollector(t *testing.T) {
	Convey("test memory collector", t, func() {
		c, err := NewMemoryCollector()
		So(err, ShouldBeNil)
		for _, m := range c.Collect() {
			fmt.Println(m)
		}
	})
}

func TestNetworkCollector(t *testing.T) {
	Convey("test network collector", t, func() {
		c, err := NewNetworkCollector("en1")
		So(err, ShouldBeNil)
		for i := 0; i < 10; i++ {
			time.Sleep(time.Second)
			for _, m := range c.Collect() {
				fmt.Println(m)
			}
		}
	})
}

func TestLoadAvgCollector(t *testing.T) {
	Convey("test loadavg collector", t, func() {
		c, err := NewLoadAvgCollector()
		So(err, ShouldBeNil)
		for _, m := range c.Collect() {
			fmt.Println(m)
		}
	})
}

func TestDiskCollector(t *testing.T) {
	Convey("test disk collector", t, func() {
		c, err := NewDiskCollector()
		So(err, ShouldBeNil)
		for _, m := range c.Collect() {
			fmt.Println(m)
		}
	})
}

func TestIOPSCollector(t *testing.T) {
	Convey("test iops collector", t, func() {
		c, err := NewIOPSCollector("sda")
		So(err, ShouldBeNil)
		for i := 0; i < 10; i++ {
			time.Sleep(time.Second)
			for _, m := range c.Collect() {
				fmt.Println(m)
			}
		}
	})
}
