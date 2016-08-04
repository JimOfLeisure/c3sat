package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/myjimnelson/c3sat/readciv3"
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
				err := ioutil.WriteFile("./out.sav", readciv3.Readciv3(c.Args().First()), 0644)
				check(err)
				log.Println("Saved to out.sav in current folder")
				return nil
			},
		},
		{
			Name:    "hexdump",
			Aliases: []string{"x"},
			Usage:   "hex dump a Civ3 data file to stdout",
			Action: func(c *cli.Context) error {
				fmt.Print(hex.Dump(readciv3.Readciv3(c.Args().First())))
				return nil
			},
		},
	}

	app.Run(os.Args)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
