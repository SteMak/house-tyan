package main

import "github.com/urfave/cli"

func commands() *cli.App {
	a := &cli.App{
		Name:        "Discord AniHouse server Bot",
		Description: "Discord AniHouse server Bot",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "d, debug",
				Usage: "Enable debug mode",
			},
			&cli.BoolFlag{
				Name:  "m,migrate",
				Usage: "Migrate database",
			},
			&cli.StringFlag{
				Name:  "c, config",
				Value: "./config/work/config.json",
				Usage: "Path to the config file",
			},
		},
		Action: run,
	}
	a.UseShortOptionHandling = true
	return a
}
