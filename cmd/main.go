package main

import (
	"context"
	"os"

	"github.com/urfave/cli/v2"
	server "github.com/ztalab/ZASentinel/internal/app"
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
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "conf",
			Aliases:  []string{"c"},
			Usage:    "The configuration file(.json,.yaml,.toml)",
			Required: true,
		},
	}
	app.Action = func(c *cli.Context) error {
		return server.Run(ctx,
			server.SetConfigFile(c.String("conf")),
			server.SetVersion(VERSION))
	}
	err := app.Run(os.Args)
	if err != nil {
		logger.WithContext(ctx).Errorf(err.Error())
	}
}
