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
			v := c.Collect()
			fmt.Println(v)
		}
	})
}

func TestMemoryCollector(t *testing.T) {
	Convey("test memory collector", t, func() {
		c, err := NewMemoryCollector()
		So(err, ShouldBeNil)
		v := c.Collect()
		fmt.Println(v)
	})
}

func TestNetworkCollector(t *testing.T) {
	Convey("test network collector", t, func() {
		c, err := NewNetworkCollector("en1")
		So(err, ShouldBeNil)
		for i := 0; i < 10; i++ {
			time.Sleep(time.Second)
			v := c.Collect()
			fmt.Println(v)
		}
	})
}

func TestLoadAvgCollector(t *testing.T) {
	Convey("test loadavg collector", t, func() {
		c, err := NewLoadAvgCollector()
		So(err, ShouldBeNil)
		v := c.Collect()
		fmt.Println(v)
	})
}

func TestDiskCollector(t *testing.T) {
	Convey("test disk collector", t, func() {
		c, err := NewDiskCollector()
		So(err, ShouldBeNil)
		v := c.Collect()
		fmt.Println(v)
	})
}
