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
		for i:=0; i < 10; i++ {
			time.Sleep(time.Second)
			So(err, ShouldBeNil)
			v := c.Collect()
			fmt.Println(v)
		}
	})
}

func TestMemoryCollector(t *testing.T) {
	Convey("test memory collector", t, func() {
		c, err := NewMemoryCollector()
		for i:=0; i < 10; i++ {
			time.Sleep(time.Second)
			So(err, ShouldBeNil)
			v := c.Collect()
			fmt.Println(v)
		}
	})
}

func TestNetworkCollector(t *testing.T) {
	Convey("test network collector", t, func() {
		c, err := NewNetworkCollector("en1")
		for i:=0; i < 10; i++ {
			time.Sleep(time.Second)
			So(err, ShouldBeNil)
			v := c.Collect()
			fmt.Println(v)
		}
	})
}
