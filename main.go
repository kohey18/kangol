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
	app.Usage = "ECS deployment tool"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "conf",
			Value: "",
			Usage: "ECS service family at task definition",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "--debug has a debug mode",
		},
	}

	app.Action = func(c *cli.Context) {
		deploy(c.String("conf"), c.Bool("debug"))
	}
	app.Run(os.Args)

	finished <- true

}
