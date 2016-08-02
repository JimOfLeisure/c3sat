// Package blast is to decompress SAV and BIQ files
// Obviously not yet complete
// Basing early steps on http://www.jonathan-petitcolas.com/2014/09/25/parsing-binary-files-in-go.html
package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	path := "/temp/civ3saves/about to win English, 1340 AD.SAV"

	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Error while opening file", err)
	}

	defer file.Close()

	fmt.Printf("%s opened\n", path)
	fmt.Printf("print me please\n")
	header := readNextBytes(file, 2)
	fmt.Printf("HaloooooOOooo %x\n", header)
}

func readNextBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}
