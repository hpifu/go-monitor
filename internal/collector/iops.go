package collector

//
//import (
//	"fmt"
//	"math"
//	"strings"
//	"time"
//
//	"github.com/mackerelio/go-osstat/disk"
//)
//
//// linux command: iostat
//type IOPSCollector struct {
//	name  string
//	value *disk.Stats
//	ts    time.Time
//
//	keys map[string]string
//}
//
//func NewIOPSCollector(name string) (*IOPSCollector, error) {
//	value := collectIOPS(name)
//	ts := time.Now()
//
//	if value == nil {
//		return nil, fmt.Errorf("network not found, name: [%v]", name)
//	}
//
//	return &IOPSCollector{
//		name:  name,
//		value: value,
//		ts:    ts,
//
//		keys: map[string]string{
//			"host": Hostname(),
//		},
//	}, nil
//}
//
//func collectIOPS(name string) *disk.Stats {
//	vals, _ := disk.Get()
//	for _, val := range vals {
//		if strings.HasPrefix(val.Name, name) {
//			return &val
//		}
//	}
//
//	return nil
//}
//
//func (c *IOPSCollector) Collect() []*Metric {
//	value := collectIOPS(c.name)
//	ts := time.Now()
//
//	res := map[string]interface{}{
//		"rps": math.Round(float64(value.ReadsCompleted-c.value.ReadsCompleted)/float64(ts.Sub(c.ts)/time.Second)*100) / 100,
//		"wps": math.Round(float64(value.WritesCompleted-c.value.WritesCompleted)/float64(ts.Sub(c.ts)/time.Second)*100) / 100,
//	}
//
//	c.value = value
//	c.ts = ts
//
//	return []*Metric{
//		{
//			Keys:      c.keys,
//			Vals:      res,
//			Timestamp: time.Now(),
//		},
//	}
//}
