package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"

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
			Name:    "dev",
			Aliases: []string{"z"},
			Usage:   "Who knows? It's whatever the dev is working on right now",
			Action: func(c *cli.Context) error {
				var gameData parseciv3.Civ3Data
				var err error
				path := c.Args().First()
				gameData, err = parseciv3.ParseCiv3(path)
				if err != nil {
					if parseErr, ok := err.(parseciv3.ParseError); ok {
						log.Printf("Expected: %s\nHex Dump:\n%s\n", parseErr.Expected, parseErr.Hexdump)
					}
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
