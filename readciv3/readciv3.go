// Package readciv3 is to decompress SAV and BIQ files
// Obviously not yet complete
package main

import (
	"bytes"
	"encoding/hex"
	"io"
	"log"
	"os"
)

func main() {
	path := "/temp/civ3saves/about to win English, 1340 AD.SAV"

	// Remove the date/time stamp from log lines
	log.SetFlags(0)

	// Open file, hanlde errors, defer close
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Error while opening file", err)
	}
	defer file.Close()
	log.Printf("%s opened\n", path)

	header := make([]byte, 2)
	_, err = file.Read(header)
	if err != nil {
		log.Fatal(err)
	}
	switch {
	case header[0] == 0x00 && header[1] == 0x06:
		log.Println("Compressed Civ3 file detected")
	case string(header) == "CI":
		log.Fatal("Uncompressed Civ3 SAV file detected")
	case string(header) == "BI":
		log.Fatal("Uncompressed Civ3 BIC file detected")
	}

	// Create bitstream reader
	civ3Bitstream := NewReader(file)
	var uncData bytes.Buffer

	for {
		foo, err := civ3Bitstream.ReadBit()
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		switch foo {
		case true:
			log.Printf("Data hex dump:\n%s\n", hex.Dump(uncData.Bytes()))
			log.Fatal("Dictionary logic not yet implemented.\n")
		case false:
			{
				aByte, err := civ3Bitstream.ReadByte()
				if err != nil {
					log.Fatal(err)
				}
				uncData.Write([]byte{aByte})
			}
		}
	}
}
