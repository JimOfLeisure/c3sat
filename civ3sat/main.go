package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"text/tabwriter"

	"github.com/myjimnelson/c3sat/parseciv3"
	"github.com/urfave/cli"
)

func main() {
	// Remove the date/time stamp from log lines
	log.SetFlags(0)

	app := cli.NewApp()
	app.Name = "Civ3 Show-And-Tell"
	app.Usage = "A utility to extract data from Civ3 SAV and BIQ files. Provide a file name of a SAV or BIQ file after the command."

	app.Commands = []cli.Command{
		{
			Name:    "seed",
			Aliases: []string{"s"},
			Usage:   "Show the world seed and map settings needed to generate the map, if this map was randomly generated.",
			Action: func(c *cli.Context) error {
				var gameData parseciv3.Civ3Data
				var err error
				path := c.Args().First()
				fmt.Println(path)
				gameData, err = parseciv3.NewCiv3Data(path)
				if err != nil {
					return err
				}
				fmt.Println()
				w := new(tabwriter.Writer)
				w.Init(os.Stdout, 0, 8, 0, '\t', 0)
				settings := gameData.WorldSettings()
				for i := range settings {
					fmt.Fprintf(w, "%s\t%s\t%s\n", settings[i][0], settings[i][1], settings[i][2])
				}
				w.Flush()
				return nil
			},
		},
		{
			Name:    "decompress",
			Aliases: []string{"d"},
			Usage:   "decompress a Civ3 data file to out.sav in the current folder",
			Action: func(c *cli.Context) error {
				filedata, _, err := parseciv3.ReadFile(c.Args().First())
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
			Action: func(c *cli.Context) error {
				filedata, _, err := parseciv3.ReadFile(c.Args().First())
				if err != nil {
					return err
				}

				fmt.Print(hex.Dump(filedata))
				return nil
			},
		},
		{
			Name:    "map",
			Aliases: []string{"m"},
			Usage:   "Dump a JSON file of map data",
			Action: func(c *cli.Context) error {
				var gameData parseciv3.Civ3Data
				var err error
				path := c.Args().First()
				fmt.Println(path)
				gameData, err = parseciv3.NewCiv3Data(path)
				if err != nil {
					if parseErr, ok := err.(parseciv3.ParseError); ok {
						return parseErr
					}
					return err
				}
				err = ioutil.WriteFile("./civmap.json", gameData.JSONMap(), 0644)
				if err != nil {
					log.Println("Error writing file")
					return err
				}

				log.Println("Saved to civmap.json in current folder")
				return nil
			},
		},
		{
			Name:    "dev",
			Aliases: []string{"z"},
			Usage:   "Who knows? It's whatever the dev is working on right now",
			Action: func(c *cli.Context) error {
				var gameData parseciv3.Civ3Data
				var err error
				path := c.Args().First()
				fmt.Println(path)
				gameData, err = parseciv3.NewCiv3Data(path)
				if err != nil {
					if parseErr, ok := err.(parseciv3.ParseError); ok {
						log.Printf("Expected: %s\nHex Dump:\n%s\n", parseErr.Expected, parseErr.Hexdump)
						return parseErr
					}
					// fmt.Print(gameData.Debug())
					return err
				}
				fmt.Print(gameData.Info())
				fmt.Print(gameData.Debug())
				return nil
			},
		},
	}

	app.Run(os.Args)
}
