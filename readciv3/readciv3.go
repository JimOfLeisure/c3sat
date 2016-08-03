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

// NOTE: Do I care what the dictionary size is? If I'm buffering the whole file?
// Answer: Yes I do care, but not because of the buffer. It's the offset bit count that matters

type lengthCode struct {
	value, extraBits int
}

var lengthLookup = map[string]lengthCode{
	"101":     lengthCode{2, 0},
	"11":      lengthCode{3, 0},
	"100":     lengthCode{4, 0},
	"011":     lengthCode{5, 0},
	"0101":    lengthCode{6, 0},
	"0100":    lengthCode{7, 0},
	"0011":    lengthCode{8, 0},
	"00101":   lengthCode{9, 0},
	"00100":   lengthCode{10, 1},
	"00011":   lengthCode{12, 2},
	"00010":   lengthCode{16, 3},
	"000011":  lengthCode{24, 4},
	"000010":  lengthCode{40, 5},
	"000001":  lengthCode{72, 6},
	"0000001": lengthCode{136, 7},
	"0000000": lengthCode{264, 8},
}

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
	case header[0] == 0x00 && (header[1] == 0x04 || header[1] == 0x05 || header[1] == 0x06):
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
			log.Printf("%v", civ3Bitstream.lengthsequence())
			log.Fatal("Dictionary logic not yet fully implemented.\n")
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

func (b *BitReader) lengthsequence() int {
	var sequence bytes.Buffer
	// TODO: Set max limit on sequence length, error out
	count := 0
	for _, zpresent := lengthLookup[sequence.String()]; !zpresent && count < 20; count++ {
		log.Printf("%v", zpresent)
		log.Println(sequence.String())
		log.Printf("%v", lengthLookup[sequence.String()])
		bit, _ := b.ReadBit()
		if bit {
			sequence.WriteString("1")
		} else {
			sequence.WriteString("0")
		}
		// hack, but not sure how to check every iteration in for params
		_, zpresent = lengthLookup[sequence.String()]
	}
	log.Printf("%v", lengthLookup["11"])
	log.Println(sequence.String())
	log.Fatalf("Decoded length sequence is %v, to read %v more bits", lengthLookup[sequence.String()].value, lengthLookup[sequence.String()].extraBits)
	return 0
}

// ReadByte reads a single byte from the stream, regardless of alignment
func (b *BitReader) ReadByte() (byte, error) {

	// If I init inside the loop these are out of scope
	var byt byte = 255
	var err error

	// Shift in 8 bits, LSBit first
	for i := 0; i < 8; i++ {
		bit, looperr := b.ReadBit()
		if looperr != nil {
			log.Fatal(looperr)
		}
		byt >>= 1
		if bit {
			byt |= 128
		}
	}

	return byt, err
}
