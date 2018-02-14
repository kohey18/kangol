package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	finished := make(chan bool)
	go loading(finished)

	app := cli.NewApp()
	app.Name = "kangol"
	app.Version = "0.1.1"
	app.Usage = "ECS deployment tool"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "conf",
			Value: "",
			Usage: "ECS service family at task definition",
		},
		cli.StringFlag{
			Name:  "tag",
			Value: "",
			Usage: "--tag has a container tag",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "--debug has a debug mode",
		},
		cli.BoolFlag{
			Name:  "skip-polling",
			Usage: "--skip-polling skip polling",
		},
		cli.BoolFlag{
			Name:  "loading",
			Usage: "--loading has a loading while deploying",
		},
	}

	app.Action = func(c *cli.Context) {

		if c.Bool("loading") == false {
			finished <- true
		}
		deploy(c.String("conf"), c.String("tag"), c.Bool("debug"), c.Bool("skip-polling"))
	}
	app.Run(os.Args)

}
