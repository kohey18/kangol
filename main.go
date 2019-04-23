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
	app.Version = "0.2.4"
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
		cli.Int64Flag{
			Name:  "polling-time",
			Value: 1,
			Usage: "--polling-time 10",
		},
	}

	app.Action = func(c *cli.Context) {

		if c.Bool("loading") == false {
			finished <- true
		}
		deploy(
			c.String("conf"),
			c.String("tag"),
			c.Bool("debug"),
			c.Bool("skip-polling"),
			c.Int64("polling-time"))
	}

	app.Commands = []cli.Command{
		{
			Name:  "run",
			Usage: "run an ECS task",
			Flags: []cli.Flag{
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
				cli.StringFlag{
					Name:  "command",
					Value: "",
					Usage: "--command has a run task override command",
				},
				cli.Int64Flag{
					Name:  "cpu",
					Value: 256,
					Usage: "--cpu has a run task override cpu reserved",
				},
				cli.Int64Flag{
					Name:  "memory",
					Value: 256,
					Usage: "--memory has a run task override memory reserved",
				},
			},
			Action: func(c *cli.Context) {
				runTask(c.String("conf"), c.String("tag"), c.String("command"), c.Int64("cpu"), c.Int64("memory"))
			},
		},
	}

	app.Run(os.Args)

}
