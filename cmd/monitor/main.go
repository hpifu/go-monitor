package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/hpifu/go-kit/logger"
	"github.com/hpifu/go-monitor/internal/collector"
	"github.com/hpifu/go-monitor/internal/monitor"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/spf13/viper"
	"github.com/yosuke-furukawa/json5/encoding/json5"
)

// AppVersion name
var AppVersion = "unknown"

type collectorInfo struct {
	Class    string        `json:"class"`
	Params   []interface{} `json:"params"`
	Table    string        `json:"table"`
	Interval string        `json:"interval"`
}

type collectorInfos struct {
	Infos []collectorInfo `json:"collector"`
}

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

	// init monitor
	m, err := monitor.NewMonitor(
		config.GetString("monitor.address"),
		config.GetString("monitor.database"),
		config.GetInt("monitor.batch"),
		config.GetInt("monitor.retry"),
	)
	if err != nil {
		panic(err)
	}
	m.SetLogger(infoLog, warnLog, accessLog)

	// use json5 because viper doesn't support json array
	var infos collectorInfos
	fp, err = os.Open(*configfile)
	if err != nil {
		panic(err)
	}
	err = json5.NewDecoder(fp).Decode(&infos)
	if err != nil {
		panic(err)
	}
	_ = fp.Close()
	for _, info := range infos.Infos {
		fmt.Println(info)
		c, err := collector.NewCollector(info.Class, info.Params)
		if err != nil {
			panic(err)
		}
		interval, err := time.ParseDuration(info.Interval)
		if err != nil {
			panic(err)
		}
		m.AddCollector(c, info.Table, interval)
	}
	if err := m.Monitor(); err != nil {
		panic(err)
	}
	infoLog.Infof("init monitor success")

	// graceful quit
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	infoLog.Infof("%v shutdown ...", os.Args[0])
	m.Stop()
	_ = warnLog.Out.(*rotatelogs.RotateLogs).Close()
	_ = accessLog.Out.(*rotatelogs.RotateLogs).Close()
	infoLog.Errorf("%v shutdown success", os.Args[0])
	_ = infoLog.Out.(*rotatelogs.RotateLogs).Close()
}
