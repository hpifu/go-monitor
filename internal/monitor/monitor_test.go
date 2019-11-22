package monitor

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMonitor(t *testing.T) {
	Convey("test monitor", t, func() {
		m, err := NewMonitor("http://localhost:8086", "mydb")
		So(err, ShouldBeNil)
		So(m.Save(&Mertic{
			Table: "testtbl",
			Keys: map[string]string{
				"field1": "key1",
			},
			Value: 12.3,
		}), ShouldBeNil)
	})
}
