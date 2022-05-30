package main

import (
	"context"
	"github.com/ztalab/ZASentinel/internal"
	"github.com/ztalab/ZASentinel/internal/client"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/ztalab/ZASentinel/pkg/logger"
)

var VERSION = "0.0.0"

func main() {
	logger.SetVersion(VERSION)
	ctx := logger.NewTagContext(context.Background(), "__main__")

	app := cli.NewApp()
	app.Name = "za-sentinel"
	app.Version = VERSION
	app.Usage = "Security, network acceleration, zero trust network architecture"
	app.Flags = commonConfig()
	app.Commands = []*cli.Command{
		client.NewCliCmd(ctx),
		newRelayCmd(ctx),
		newServerCmd(ctx),
	}
	err := app.Run(os.Args)
	if err != nil {
		logger.WithContext(ctx).Errorf("%v", err)
	}
}

func newRelayCmd(ctx context.Context) *cli.Command {
	return &cli.Command{
		Name:  "relay",
		Usage: "Run relay server",
		Action: func(c *cli.Context) error {
			return internal.Run(ctx,
				internal.SetConfigFile(c.String("conf")),
				internal.SetVersion(VERSION))
		},
	}
}

func newServerCmd(ctx context.Context) *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Run server server",
		Action: func(c *cli.Context) error {
			return internal.Run(ctx,
				internal.SetConfigFile(c.String("conf")),
				internal.SetVersion(VERSION))
		},
	}
}

func commonConfig() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "conf",
			Aliases:  []string{"c"},
			Usage:    "App configuration file(.json,.yaml,.toml)",
			Required: true,
		},
	}
}
