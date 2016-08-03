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

var offsetLookup = map[string]int{
	"11":       0x00,
	"1011":     0x01,
	"1010":     0x02,
	"10011":    0x03,
	"10010":    0x04,
	"10001":    0x05,
	"10000":    0x06,
	"011111":   0x07,
	"011110":   0x08,
	"011101":   0x09,
	"011100":   0x0a,
	"011011":   0x0b,
	"011010":   0x0c,
	"011001":   0x0d,
	"011000":   0x0e,
	"010111":   0x0f,
	"010110":   0x10,
	"010101":   0x11,
	"010100":   0x12,
	"010011":   0x13,
	"010010":   0x14,
	"010001":   0x15,
	"0100001":  0x16,
	"0100000":  0x17,
	"0011111":  0x18,
	"0011110":  0x19,
	"0011101":  0x1a,
	"0011100":  0x1b,
	"0011011":  0x1c,
	"0011010":  0x1d,
	"0011001":  0x1e,
	"0011000":  0x1f,
	"0010111":  0x20,
	"0010110":  0x21,
	"0010101":  0x22,
	"0010100":  0x23,
	"0010011":  0x24,
	"0010010":  0x25,
	"0010001":  0x26,
	"0010000":  0x27,
	"0001111":  0x28,
	"0001110":  0x29,
	"0001101":  0x2a,
	"0001100":  0x2b,
	"0001011":  0x2c,
	"0001010":  0x2d,
	"0001001":  0x2e,
	"0001000":  0x2f,
	"00001111": 0x30,
	"00001110": 0x31,
	"00001101": 0x32,
	"00001100": 0x33,
	"00001011": 0x34,
	"00001010": 0x35,
	"00001001": 0x36,
	"00001000": 0x37,
	"00000111": 0x38,
	"00000110": 0x39,
	"00000101": 0x3a,
	"00000100": 0x3b,
	"00000011": 0x3c,
	"00000010": 0x3d,
	"00000001": 0x3e,
	"00000000": 0x3f,
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
			length := civ3Bitstream.lengthsequence()
			// log.Printf("Length %v", length)
			if length == 519 {
				log.Printf("Data hex dump:\n%s\n", hex.Dump(uncData.Bytes()))
				log.Fatal("End of stream token reached")
			}
			_ = civ3Bitstream.offsetsequence(int(header[1]))
			// offset := civ3Bitstream.offsetsequence(int(header[1]))
			// log.Printf("Offset %v", offset)
			// log.Fatal("Dictionary logic not yet fully implemented.\n")
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
	// TODO: Do I care about err handling? Currently using _
	count := 0
	for _, zpresent := lengthLookup[sequence.String()]; !zpresent && count < 8; count++ {
		bit, _ := b.ReadBit()
		if bit {
			sequence.WriteString("1")
		} else {
			sequence.WriteString("0")
		}
		// hack, but not sure how to check every iteration in for params
		_, zpresent = lengthLookup[sequence.String()]
	}
	xxxes, _ := b.ReadBits(uint(lengthLookup[sequence.String()].extraBits))
	// log.Printf("Decoded length sequence is %v, to read %v more bits", lengthLookup[sequence.String()].value, lengthLookup[sequence.String()].extraBits)
	return lengthLookup[sequence.String()].value + int(xxxes)
}
func (b *BitReader) offsetsequence(dictsize int) int {
	var sequence bytes.Buffer
	// TODO: Do I care about err handling? Currently using _
	count := 0
	for _, zpresent := offsetLookup[sequence.String()]; !zpresent && count < 8; count++ {
		bit, _ := b.ReadBit()
		if bit {
			sequence.WriteString("1")
		} else {
			sequence.WriteString("0")
		}
		// hack, but not sure how to check every iteration in for params
		_, zpresent = offsetLookup[sequence.String()]
	}
	loworderbits, _ := b.ReadBits(uint(dictsize))
	// log.Printf("Decoded length sequence is %v, to read %v more bits", offsetLookup[sequence.String()].value, offsetLookup[sequence.String()].extraBits)
	return offsetLookup[sequence.String()]<<uint(dictsize) + int(loworderbits)
}

// ReadByte reads a single byte from the stream, regardless of alignment
func (b *BitReader) ReadByte() (byte, error) {

	// If I init inside the loop these are out of scope
	var byt byte
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

// ReadBits reads  nbits from the stream
func (b *BitReader) ReadBits(nbits uint) (uint, error) {

	// If I init inside the loop these are out of scope
	var value uint
	var err error

	// Use power to assign bits lsb to msb
	for i := uint(0); i < nbits; i++ {
		bit, _ := b.ReadBit()
		if bit {
			value |= 1 << i
		}
	}

	return value, err
}
