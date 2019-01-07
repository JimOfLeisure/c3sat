package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"text/tabwriter"

	"github.com/myjimnelson/c3sat/civ3satgql"
	"github.com/myjimnelson/c3sat/parseciv3"
	"github.com/urfave/cli"
)

var saveFilePath string

var pathFlag = cli.StringFlag{
	Name: "path, p",
	// Value:       ".",
	Usage:       "`FILEPATH` of save",
	EnvVar:      "CIV3SAT_SAV",
	Destination: &saveFilePath,
}

func main() {
	app := cli.NewApp()
	app.Name = "Civ3 Show-And-Tell"
	app.Version = "0.4.0"
	app.Usage = "A utility to extract data from Civ3 SAV and BIQ files. Provide a file name of a SAV or BIQ file after the command."

	app.Commands = []cli.Command{
		{
			Name:    "seed",
			Aliases: []string{"s"},
			Usage:   "Show the world seed and map settings needed to generate the map, if was randomly generated.",
			Flags: []cli.Flag{
				pathFlag,
			},
			Action: func(c *cli.Context) error {
				fmt.Println()
				w := new(tabwriter.Writer)
				defer w.Flush()
				w.Init(os.Stdout, 0, 8, 0, '\t', 0)
				settings, err := civ3satgql.WorldSettings(saveFilePath)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				for i := range settings {
					fmt.Fprintf(w, "%s\t%s\t%s\n", settings[i][0], settings[i][1], settings[i][2])
				}
				return nil
			},
		},
		{
			Name:    "decompress",
			Aliases: []string{"d"},
			Usage:   "decompress a Civ3 data file to out.sav in the current folder",
			Flags: []cli.Flag{
				pathFlag,
			},
			Action: func(c *cli.Context) error {
				filedata, _, err := parseciv3.ReadFile(saveFilePath)
				if err != nil {
					return err
				}

				err = ioutil.WriteFile("./out.sav", filedata, 0644)
				if err != nil {
					log.Println("Error writing file")
					return err
				}

				log.Println("Saved to out.sav in current folder")
				return nil
			},
		},
		{
			Name:    "hexdump",
			Aliases: []string{"x"},
			Usage:   "hex dump a Civ3 data file to stdout",
			Flags: []cli.Flag{
				pathFlag,
			},
			Action: func(c *cli.Context) error {
				filedata, _, err := parseciv3.ReadFile(saveFilePath)
				if err != nil {
					return err
				}

				fmt.Print(hex.Dump(filedata))
				return nil
			},
		},
		{
			Name:      "graphql",
			Aliases:   []string{"gql", "g"},
			ArgsUsage: "<query>",
			Usage:     "Execute GraphQL query",
			Flags: []cli.Flag{
				pathFlag,
			},
			Action: func(c *cli.Context) error {
				// var gameData parseciv3.Civ3Data
				var err error
				query := c.Args().First()
				result, err := civ3satgql.Query(query, saveFilePath)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				fmt.Print(result)
				return nil
			},
		},
		{
			Name:    "api",
			Aliases: []string{"www"},
			Usage:   "Open save, start GraphQL API at http://127.0.0.1:8080/graphql . Control-c to exit.",
			Flags: []cli.Flag{
				pathFlag,
				cli.StringFlag{
					Name:   "addr",
					Value:  "127.0.0.1",
					Usage:  "`ADDRESS` on which to bind",
					EnvVar: "CIV3SAT_ADDR",
				},
				cli.StringFlag{
					Name:   "port",
					Value:  "8080",
					Usage:  "`PORT` on which to listen",
					EnvVar: "CIV3SAT_PORT",
				},
			},
			Action: func(c *cli.Context) error {
				var err error
				fmt.Println("Starting API server for save file at " + saveFilePath)
				fmt.Println("GraphQL at http://" + c.String("addr") + ":" + c.String("port") + "/graphql")
				fmt.Println("Press control-C to exit")
				err = civ3satgql.Server(saveFilePath, c.String("addr"), c.String("port"))
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				return nil
			},
		},
	}
	app.Run(os.Args)
}
