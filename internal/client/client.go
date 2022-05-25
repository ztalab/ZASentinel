package client

import (
	"context"
	"github.com/urfave/cli/v2"
	"github.com/ztalab/ZASentinel/pkg/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func NewCliCmd(ctx context.Context) *cli.Command {
	return &cli.Command{
		Name:  "cli",
		Usage: "Run cli server",
		Subcommands: []*cli.Command{
			UpCmd(ctx),
		},
	}
}

func Run(ctx context.Context, f func(ctx context.Context) (func(), error)) error {
	state := 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	cleanFunc, err := f(ctx)
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
