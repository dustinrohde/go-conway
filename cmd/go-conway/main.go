package main

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/dustinrohde/go-conway"
	util "github.com/dustinrohde/go-conway/util/go-conway"
	"github.com/urfave/cli"
)

func main() {
	app := initApp()
	app.Run(os.Args)
}

func initApp() *cli.App {
	app := cli.NewApp()
	app.Name = "conway"
	app.Usage = "Run Conway's Game of Life simulation."

	var grid conway.Grid
	config := conway.DefaultRunConfig()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "outfile, o",
			Value: "-",
			Usage: "File to write results to. If `FILE` is `-`, use stdout.",
		},
		cli.DurationFlag{
			Name:        "delay, d",
			Value:       config.Delay,
			Usage:       "`TIME` to pause between turns; e.g. 500ms, 3s, 1m",
			Destination: &config.Delay,
		},
		cli.IntFlag{
			Name:        "turns, t",
			Value:       config.MaxTurns,
			Usage:       "Max `TURNS` to run. If < 0, run indefinitely.",
			Destination: &config.MaxTurns,
		},
		cli.BoolFlag{
			Name:        "clear, c",
			Usage:       "Clear screen between each turn.",
			Destination: &config.ClearScreen,
		},
		cli.BoolFlag{
			Name:        "interactive, i",
			Usage:       "Wait for input between each turn.",
			Destination: &config.Interactive,
		},
		cli.BoolFlag{
			Name:        "spinner, s",
			Usage:       "Show an animated spinner between turns.",
			Destination: &config.Spinner,
		},
	}

	app.Commands = []cli.Command{
		cli.Command{
			Name:    "defined",
			Aliases: []string{"def"},
			Usage:   "start with a predefined grid",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "grid, g",
					Value: "",
					Usage: strings.Join([]string{
						"Starting grid.",
						"\t\tIf `FILE` starts with '@', interpret it as a file path.",
						"\t\tIf `FILE` is '-', read from stdout.",
						"\t\tIf `FILE` is absent or blank, use a demo starting grid.",
					}, "\n"),
				},
			},
			Action: func(c *cli.Context) error {
				if path := c.String("grid"); len(path) > 0 {
					if path == "-" {
						// Read Grid from stdin.
						config.GridFile = os.Stdin
					} else if path[0] == '@' {
						// Read Grid from file.
						file, err := os.OpenFile(path[1:], os.O_RDONLY, 0600)
						util.Guard(err)
						defer file.Close()
						config.GridFile = file
					} else {
						// Read Grid from argument.
						strings.Split
					}
				} else {
					// Use the demo Grid.
				}

				gridBytes, err := ioutil.ReadAll(config.GridFile)
				util.Guard(err)
				grid = conway.FromString(string(gridBytes))
			},
		},
		cli.Command{
			Name:    "random",
			Aliases: []string{"rand"},
			Usage:   "start with a randomly generated grid",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "width, w",
					Value: 9,
					Usage: "max `WIDTH` of the grid",
				},
				cli.IntFlag{
					Name:  "height, h",
					Value: 9,
					Usage: "max `HEIGHT` of the grid",
				},
				cli.Float64Flag{
					Name:  "probability, p",
					Value: 0.5,
					Usage: "probability of living cells, where 0 < `PROB` <= 1",
				},
			},
			Action: func(c *cli.Context) error {
				grid = RandomGrid(c.Int("width"), c.Int("height"), c.Float64("probability"))
			},
		},
	}

	app.Action = func(c *cli.Context) error {
		if path := c.String("outfile"); path != "-" {
			// Write results to a file.
			file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			util.Guard(err)
			config.OutFile = file
			defer file.Close()
		}

		conway.Run(grid, config)
		return nil
	}

	return app
}
