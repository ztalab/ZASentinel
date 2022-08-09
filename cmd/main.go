// Copyright 2022-present The Ztalab Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
