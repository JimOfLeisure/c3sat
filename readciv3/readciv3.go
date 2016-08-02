// Package readciv3 is to decompress SAV and BIQ files
// Obviously not yet complete
package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	path := "/temp/civ3saves/about to win English, 1340 AD.SAV"

	// Open file, hanlde errors, defer close
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Error while opening file", err)
	}
	defer file.Close()
	fmt.Printf("%s opened\n", path)

	header := make([]byte, 2)
	_, err = file.Read(header)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%x\n", header)

	// Create bitstream reader
	civ3Bitstream := NewReader(file)

	for {
		foo, err := civ3Bitstream.ReadBit()
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		// fmt.Printf("%v %v\n", foo, err)
		switch foo {
		case true:
			fmt.Println("\n\n")
			log.Fatal("Dictionary logic not yet implemented.\n")
		case false:
			{
				aByte, err := civ3Bitstream.ReadByte()
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("%02x ", aByte)
			}
		}
	}
}
