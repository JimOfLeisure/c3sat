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
	// Remove the date/time stamp from log lines
	log.SetFlags(0)

	log.Println("Early test program. Call this program with a Civ 3 SAV file or BIC/X/Q file.")
	log.Println("It will report if the file is compressed or uncompressed,")
	log.Println("and print the first several bytes of any compressed files in hex dump format.")
	path := os.Args[1]

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
	default:
		log.Fatalf("Unrecognized file type. Two byte header is %v", header)
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
