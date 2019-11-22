module github.com/hpifu/go-monitor

go 1.12

replace (
	cloud.google.com/go => github.com/googleapis/google-cloud-go v0.0.0-20190603211518-c8433c9aaceb
	go.etcd.io/bbolt => github.com/etcd-io/bbolt v1.3.4-0.20191001164932-6e135e5d7e3d
	go.uber.org/atomic => github.com/uber-go/atomic v1.4.1-0.20190731194737-ef0d20d85b01
	go.uber.org/multierr => github.com/uber-go/multierr v1.2.0
	go.uber.org/zap => github.com/uber-go/zap v1.10.1-0.20190926184545-d8445f34b4ae
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190605123033-f99c8df09eb5
	golang.org/x/exp => github.com/golang/exp v0.0.0-20190510132918-efd6b22b2522
	golang.org/x/image => github.com/golang/image v0.0.0-20190523035834-f03afa92d3ff
	golang.org/x/lint => github.com/golang/lint v0.0.0-20190409202823-959b441ac422
	golang.org/x/mobile => github.com/golang/mobile v0.0.0-20190607214518-6fa95d984e88
	golang.org/x/net => github.com/golang/net v0.0.0-20190606173856-1492cefac77f
	golang.org/x/oauth2 => github.com/golang/oauth2 v0.0.0-20190604053449-0f29369cfe45
	golang.org/x/sync => github.com/golang/sync v0.0.0-20190423024810-112230192c58
	golang.org/x/sys => github.com/golang/sys v0.0.0-20190602015325-4c4f7f33c9ed
	golang.org/x/text => github.com/golang/text v0.3.2
	golang.org/x/time => github.com/golang/time v0.0.0-20190308202827-9d24e82272b4
	golang.org/x/tools => github.com/golang/tools v0.0.0-20190608022120-eacb66d2a7c3
	google.golang.org/api => github.com/googleapis/google-api-go-client v0.6.0
	google.golang.org/appengine => github.com/golang/appengine v1.6.1
	google.golang.org/genproto => github.com/google/go-genproto v0.0.0-20190605220351-eb0b1bdb6ae6
	google.golang.org/grpc => github.com/grpc/grpc-go v1.21.1
)

require (
	github.com/influxdata/influxdb1-client v0.0.0-20190809212627-fc22c7df067e
	github.com/mackerelio/go-osstat v0.0.0-20191011055430-af968620273d
	github.com/smartystreets/goconvey v1.6.4
)
