package main

import (
	"fmt"
	"os"

	"airman.io/lite/lite"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "lite"
	app.Version = "0.0.1.SNAPSHOT"
	app.Usage = "make an explosive entrance"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Load configuration from `FILE`",
			Value: "config.toml",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "watch",
			Aliases: []string{"w"},
			Usage:   "Watch the choosen environment",
			Action:  lite.WatchAction,
			Before: func(c *cli.Context) error {

				configFilePath := c.GlobalString("config")
				err := lite.Load(configFilePath)
				if err != nil {
					return cli.NewExitError("can not find the config.yaml", 1)
				}
				return nil
			},
		},
		{
			Name:    "build",
			Aliases: []string{"b"},
			Usage:   "Build the choosen environment",
			Action: func(c *cli.Context) error {
				fmt.Println("build command")
				return nil
			},
		},
	}

	app.Action = func(c *cli.Context) error {
		fmt.Println(c.String("config"))
		cli.ShowAppHelp(c)
		return nil
	}

	app.Run(os.Args)
}
