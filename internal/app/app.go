package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/ztalab/ZASentinel/internal/app/bll"
	"github.com/ztalab/ZASentinel/internal/app/config"
	"github.com/ztalab/ZASentinel/internal/app/initer"
	"github.com/ztalab/ZASentinel/pkg/influxdb"
	influx_client "github.com/ztalab/ZASentinel/pkg/influxdb/client/v2"
	"github.com/ztalab/ZASentinel/pkg/logger"
	"github.com/ztalab/ZASentinel/pkg/util/structure"
)

type options struct {
	ConfigFile string
	ModelFile  string
	Version    string
}

// Option Defining configuration items
type Option func(*options)

// SetConfigFile setting the configuration file
func SetConfigFile(s string) Option {
	return func(o *options) {
		o.ConfigFile = s
	}
}

// SetVersion set version number
func SetVersion(s string) Option {
	return func(o *options) {
		o.Version = s
	}
}

// Init application initialization
func Init(ctx context.Context, opts ...Option) (func(), error) {
	var o options
	for _, opt := range opts {
		opt(&o)
	}
	config.MustLoad(o.ConfigFile)
	// working with environment variables
	err := config.ParseConfigByEnv()
	if err != nil {
		return nil, err
	}
	config.PrintWithJSON()
	logger.WithContext(ctx).Printf("Service started, running mode：%s，version：%s，process number：%d", config.C.RunMode, o.Version, os.Getpid())

	// initialize the log module
	loggerCleanFunc, err := InitLogger()
	if err != nil {
		return nil, err
	}

	// initialize the timing module
	metrice, influxdbCleanFunc, err := InitInfluxdb(ctx)
	if err != nil {
		return nil, err
	}
	config.Is.Metrics = metrice

	basicConf, attr, err := initer.InitCert([]byte(config.C.Certificate.CertPem))
	if err != nil {
		return nil, err
	}
	var serverCleanFunc func()
	switch basicConf.Type {
	case initer.TypeClient:
		serverCleanFunc = bll.NewClient().Listen(ctx, attr)
		fmt.Println("########## start the client proxy #########")
	case initer.TypeServer:
		serverCleanFunc = bll.NewServer().Listen(ctx, attr)
		fmt.Println("########## start the server proxy #########")
	case initer.TypeRelay:
		fmt.Println("########## start the relay proxy #########")
		serverCleanFunc = bll.NewRelay().Listen(ctx, attr)
	}
	return func() {
		serverCleanFunc()
		loggerCleanFunc()
		influxdbCleanFunc()
	}, nil
}

func InitInfluxdb(ctx context.Context) (*influxdb.Metrics, func(), error) {
	if !config.C.Influxdb.Enabled {
		logger.WithContext(ctx).Warn("Influxdb Function is disabled")
		return nil, func() {}, nil
	}
	client, err := influx_client.NewHTTPClient(influx_client.HTTPConfig{
		Addr:                fmt.Sprintf("http://%v:%v", config.C.Influxdb.Address, config.C.Influxdb.Port),
		Username:            config.C.Influxdb.Username,
		Password:            config.C.Influxdb.Password,
		MaxIdleConns:        config.C.Influxdb.MaxIdleConns,
		MaxIdleConnsPerHost: config.C.Influxdb.MaxIdleConns,
	})
	if err != nil {
		return nil, func() {}, err
	}
	if _, _, err := client.Ping(1 * time.Second); err != nil {
		_ = client.Close()
		return nil, func() {}, err
	}
	iconfig := new(influxdb.CustomConfig)
	structure.Copy(config.C.Influxdb, iconfig)
	metrics, err := influxdb.NewMetrics(&influxdb.HTTPClient{
		Client: client,
		BatchPointsConfig: influx_client.BatchPointsConfig{
			Precision: config.C.Influxdb.Precision,
			Database:  config.C.Influxdb.Database,
		},
	}, iconfig)
	return metrics, func() {
		client.Close()
	}, err
}

func Run(ctx context.Context, opts ...Option) error {
	state := 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	cleanFunc, err := Init(ctx, opts...)
	if err != nil {
		return err
	}

EXIT:
	for {
		sig := <-sc
		logger.WithContext(ctx).Infof("received signal[%s]", sig.String())
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			state = 0
			break EXIT
		case syscall.SIGHUP:
		default:
			break EXIT
		}
	}

	cleanFunc()
	logger.WithContext(ctx).Infof("shutdown!")
	time.Sleep(time.Second)
	os.Exit(state)
	return nil
}
