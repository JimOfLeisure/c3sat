// Package civ3decompress is to decompress SAV and BIQ files from the game Civilization III
package civ3decompress

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
)

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

// Decompress is implemented based on the description of PKWare Data Compression Library at https://groups.google.com/forum/#!msg/comp.compression/M5P064or93o/W1ca1-ad6kgJ
// However this is only a partial implementation; The Huffman-coded literals of header 0x01 are not implemented here as they are not needed for my purpose
func Decompress(path string) []byte {
	// Open file, hanlde errors, defer close
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Error while opening file", err)
	}
	defer file.Close()
	log.Printf("File %s\n", path)

	header := make([]byte, 2)
	_, err = file.Read(header)
	if err != nil {
		log.Fatal(err)
	}
	switch {
	case header[0] == 0x00 && (header[1] == 0x04 || header[1] == 0x05 || header[1] == 0x06):
		// log.Println("Compressed file detected")
	default:
		// log.Println("Not a compressed file. Proceeding with uncompressed stream.")
		// TODO: I'm sure I'm doing this in a terribly inefficient way. Need to refactor everything to pass around file pointers I think
		theseBytes, err := ioutil.ReadFile(path)
		check(err)
		return theseBytes
	}

	// Create bitstream reader
	civ3Bitstream := NewReader(file)
	// Output bytes buffer
	var uncData bytes.Buffer
	// define length here to use in for loop setup
	var length int

	// The token equating to length 519 is the end-of-stream token
	for length != 519 {
		foo, err := civ3Bitstream.ReadBit()
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		switch foo {
		// bit 1 indicates length/offset sequences follow
		case true:
			length = civ3Bitstream.lengthsequence()
			// log.Printf("length %v", length)
			// The token equating to length 519 is the end-of-stream token
			if length != 519 {
				// If length is 2, then only two low-order bits are read for offset
				dictsize := header[1]
				if length == 2 {
					dictsize = 2
				}
				offset := civ3Bitstream.offsetsequence(int(dictsize))
				// log.Printf("offset %v", offset)
				for i := 0; i < length; i++ {
					// dictionary is just a reader for the output buffer.
					// since using .Bytes(), have to do this every loop...surely there is better way
					dict := bytes.NewReader(uncData.Bytes())
					// Position dictionary/buffer reader. 2 means from end of buffer/stream
					// offset 0 is last byte, so using -1 -offset to position for last byte
					dict.Seek(int64(-1-offset), 2)
					byt, err := dict.ReadByte()
					if err != nil {
						log.Fatal(err)
					}
					// Wouldn't think this is necessary, but let's try
					dict.Seek(int64(0), 2)
					uncData.WriteByte(byt)
					// log.Printf("byt %v", byt)
				}
			}
		// bit 0 inticates next 8 bits are literal byte, lsb first
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
	// log.Printf("Data hex dump:\n%s\n", hex.Dump(uncData.Bytes()))
	// err = ioutil.WriteFile("./out.sav", uncData.Bytes(), 0644)
	// check(err)
	return uncData.Bytes()

}

func (b *BitReader) lengthsequence() int {
	var sequence bytes.Buffer
	// TODO: Do I care about err handling? Currently using _
	count := 0
	for _, keyPresent := lengthLookup[sequence.String()]; !keyPresent; count++ {
		bit, _ := b.ReadBit()
		if bit {
			sequence.WriteString("1")
		} else {
			sequence.WriteString("0")
		}
		// hack, but not sure how to check every iteration in for params
		_, keyPresent = lengthLookup[sequence.String()]
		if count > 8 {
			log.Fatal("Did not match offset sequence")
		}
	}
	xxxes, _ := b.ReadBits(uint(lengthLookup[sequence.String()].extraBits))
	// log.Printf("Decoded length sequence is %v, to read %v more bits which are %v", lengthLookup[sequence.String()].value, lengthLookup[sequence.String()].extraBits, xxxes)
	return lengthLookup[sequence.String()].value + int(xxxes)
}
func (b *BitReader) offsetsequence(dictsize int) int {
	var sequence bytes.Buffer
	// TODO: Do I care about err handling? Currently using _
	count := 0
	for _, keyPresent := offsetLookup[sequence.String()]; !keyPresent; count++ {
		bit, _ := b.ReadBit()
		if bit {
			sequence.WriteString("1")
		} else {
			sequence.WriteString("0")
		}
		// hack, but not sure how to check every iteration in for params
		_, keyPresent = offsetLookup[sequence.String()]
		if count > 8 {
			log.Fatal("Did not match offset sequence")
		}
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

func check(e error) {
	if e != nil {
		panic(e)
	}
}
