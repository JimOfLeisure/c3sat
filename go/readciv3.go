// Package readciv3 is to decompress SAV and BIQ files
// Obviously not yet complete
package main

import (
	"fmt"
	"log"
	"os"
	// "./lsbitstream.go"
	// bitstream "github.com/dgryski/go-bitstream"
)

func main() {
	path := "/temp/civ3saves/about to win English, 1340 AD.SAV"

	file, err := os.Open(path)
	myBitstream := lsbitstream.NewReader(file)
	if err != nil {
		log.Fatal("Error while opening file", err)
	}

	defer file.Close()

	fmt.Printf("%s opened\n", path)
	// header := readNextBytes(file, 2)
	header, err := myBitstream.ReadByte()
	fmt.Printf("%x\n", header)
	header, err = myBitstream.ReadByte()
	fmt.Printf("%x\n", header)
	// myTest := bitstream.NewReader(strings.NewReader("Hi"))
	// grrr := strings.NewReader("Hi")
	foo, bar := myBitstream.ReadBit()
	fmt.Printf("%v %v\n", foo, bar)
}

func readNextBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}
