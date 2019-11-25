package main

import (
	"flag"
	"fmt"
	"github.com/hpifu/go-monitor/internal/reporter"
	"github.com/hpifu/go-monitor/internal/scheduler"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/hpifu/go-kit/logger"
	"github.com/hpifu/go-monitor/internal/collector"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/spf13/viper"
)

// AppVersion name
var AppVersion = "unknown"

type collectorInfo struct {
	Class    string        `json:"class"`
	Params   []interface{} `json:"params"`
	Table    string        `json:"table"`
	Interval time.Duration `json:"interval"`
	Enable   bool          `json:"enable"`
}

//type collectorInfos struct {
//	Infos []collectorInfo `json:"collector"`
//}

func main() {
	version := flag.Bool("v", false, "print current version")
	configfile := flag.String("c", "configs/monitor.json", "config file path")
	flag.Parse()
	if *version {
		fmt.Println(AppVersion)
		os.Exit(0)
	}

	// load config
	config := viper.New()
	config.SetEnvPrefix("account")
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	config.AutomaticEnv()
	config.SetConfigType("json")
	fp, err := os.Open(*configfile)
	if err != nil {
		panic(err)
	}
	err = config.ReadConfig(fp)
	if err != nil {
		panic(err)
	}
	_ = fp.Close()

	// init logger
	infoLog, warnLog, accessLog, err := logger.NewLoggerGroupWithViper(config.Sub("logger"))
	if err != nil {
		panic(err)
	}

	// init scheduler
	s := scheduler.NewScheduler(config.GetInt("scheduler.batch"))
	s.SetLogger(infoLog, warnLog, accessLog)

	// set reporter
	r, err := reporter.NewInfluxdbReporter(
		config.GetString("reporter.address"),
		config.GetString("reporter.database"),
		config.GetInt("reporter.retry"),
	)
	if err != nil {
		panic(err)
	}

	s.SetReporter(r)

	for key := range config.Sub("collector").AllSettings() {
		info := &collectorInfo{}
		if err := config.Sub("collector." + key).Unmarshal(info); err != nil {
			panic(err)
		}
		if !info.Enable {
			continue
		}
		fmt.Println(info)
		c, err := collector.NewCollector(info.Class, info.Params)
		if err != nil {
			panic(err)
		}
		s.AddCollector(c, info.Table, info.Interval)
	}
	infoLog.Infof("init scheduler success")

	// schedule
	if err := s.Scheduler(); err != nil {
		panic(err)
	}
	infoLog.Infof("start scheduler success")

	// graceful quit
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	infoLog.Infof("%v shutdown ...", os.Args[0])
	s.Stop()
	_ = warnLog.Out.(*rotatelogs.RotateLogs).Close()
	_ = accessLog.Out.(*rotatelogs.RotateLogs).Close()
	infoLog.Errorf("%v shutdown success", os.Args[0])
	_ = infoLog.Out.(*rotatelogs.RotateLogs).Close()
}
