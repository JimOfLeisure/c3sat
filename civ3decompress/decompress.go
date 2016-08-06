// Package civ3decompress is to decompress SAV and BIQ files from the game Civilization III
package civ3decompress

import (
	"bytes"
	"io"
	"log"
)

// Decompress is implemented based on the description of PKWare Data Compression Library at https://groups.google.com/forum/#!msg/comp.compression/M5P064or93o/W1ca1-ad6kgJ
// However this is only a partial implementation; The Huffman-coded literals of header 0x01 are not implemented here as they are not needed for my purpose
func Decompress(file io.Reader) ([]byte, error) {

	// Create bitstream reader
	civ3Bitstream := NewReader(file)
	// Output bytes buffer
	var uncData bytes.Buffer
	// define length here to use in for loop setup
	var length int

	header := make([]byte, 2)
	_, err := file.Read(header)
	if err != nil {
		return nil, FileError{err}
	}

	// The token equating to length 519 is the end-of-stream token
	const lengthEndOfStream = 519
	for length != lengthEndOfStream {
		tokenFlag, err := civ3Bitstream.ReadBit()
		if err != nil {
			return nil, FileError{err}
		}
		switch tokenFlag {
		// bit 1 indicates length/offset sequences follow
		case true:
			length = civ3Bitstream.lengthsequence()
			// log.Printf("length %v", length)
			// The token equating to length 519 is the end-of-stream token
			if length != lengthEndOfStream {
				// If length is 2, then only two low-order bits are read for offset
				dictsize := header[1]
				if length == 2 {
					dictsize = 2
				}
				offset := civ3Bitstream.offsetsequence(int(dictsize))
				for i := 0; i < length; i++ {
					// dictionary is just a reader for the output buffer.
					// since using .Bytes(), have to do this every loop...surely there is better way
					// uncData bytes.Buffer does not have Seek function
					dict := bytes.NewReader(uncData.Bytes())
					// Position dictionary/buffer reader. 2 means from end of buffer/stream
					// offset 0 is last byte, so using -1 -offset to position for last byte
					dict.Seek(int64(-1-offset), 2)
					byt, err := dict.ReadByte()
					if err != nil {
						return nil, FileError{err}
					}
					uncData.WriteByte(byt)
				}
			}
		// bit 0 inticates next 8 bits are literal byte, lsb first
		case false:
			{
				literalByte, err := civ3Bitstream.ReadByte()
				if err != nil {
					return nil, FileError{err}
				}
				uncData.Write([]byte{literalByte})
			}
		}
	}
	// log.Printf("Data hex dump:\n%s\n", hex.Dump(uncData.Bytes()))
	// err = ioutil.WriteFile("./out.sav", uncData.Bytes(), 0644)
	// check(err)
	return uncData.Bytes(), error(nil)

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
