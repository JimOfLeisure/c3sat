// Package civ3decompress is to decompress SAV and BIQ files from the game Civilization III
package civ3decompress

// Copyright (c) 2016, 2020 Jim Nelson

// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:

// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

import (
	"bytes"
	"io"
)

// DecompressByteArray is a May 2020 attempt at speeding up decompression by eliminating the reader interfaces
func DecompressByteArray(bytes []byte) ([]byte, error) {
	// The token equating to length 519 is the end-of-stream token
	const lengthEndOfStream = 519
	var outBytes []byte
	var length, bitOff, byteOff int
	var tokenFlag, bit bool
	var dictsize byte
	if len(bytes) < 3 || !(bytes[0] == 0x00 && (bytes[1] == 0x04 || bytes[1] == 0x05 || bytes[1] == 0x06)) {
		return nil, DecodeError{"Not a valid compressed byte array"}
	}
	byteOff = 2
	for {
		tokenFlag = bytes[byteOff] & (1 << bitOff) != 0
		bitOff++
		if bitOff == 8 {
			bitOff = 0
			byteOff++
		}
		if tokenFlag {
			// read length
			var sequence bitKey
			var lc lengthCode
			var keyPresent bool
			for {
				bit = bytes[byteOff] & (1 << bitOff) != 0
				sequence.key = sequence.key << 1
				sequence.keyBitLength++
				if bit {
					sequence.key |= 1
				}
				bitOff++
				if bitOff == 8 {
					bitOff = 0
					byteOff++
				}
				if sequence.keyBitLength > 8 {
					return nil, DecodeError{"Token did not match length sequence"}
				}
				if lc, keyPresent = lengthLookup[sequence]; keyPresent {
					break
				}
			}
			var xxxes byte
			for i := 0; i < lc.extraBits; i++ {
				bit = bytes[byteOff] & (1 << bitOff) != 0
				xxxes = xxxes >> 1
				if bit {
					xxxes |= 0x80
				}
				bitOff++
				if bitOff == 8 {
					bitOff = 0
					byteOff++
				}
			}
			xxxes = xxxes >> (8 - lc.extraBits)
			length = lc.value + int(xxxes)
			if length == lengthEndOfStream {
				break
			}
			// If length is 2, then only two low-order bits are read for offset
			if length == 2 {
				dictsize = 2
			} else {
				dictsize = bytes[1]
			}
			
			sequence = bitKey{0,0}
			var offset int
			for {
				bit = bytes[byteOff] & (1 << bitOff) != 0
				sequence.key = sequence.key << 1
				sequence.keyBitLength++
				if bit {
					sequence.key = sequence.key | 0b1
				}
				bitOff++
				if bitOff == 8 {
					bitOff = 0
					byteOff++
				}
				if sequence.keyBitLength > 8 {
					return nil, DecodeError{"Token did not match offset sequence"}
				}
				if offset, keyPresent = offsetLookup[sequence]; keyPresent {
					break
				}
			}
			var loworderbits uint32
			for i:=0; i < int(dictsize); i++ {
				bit = bytes[byteOff] & (1 << bitOff) != 0
				loworderbits = loworderbits >> 1
				if bit {
					loworderbits |= 0x80000000
				}
				bitOff++
				if bitOff == 8 {
					bitOff = 0
					byteOff++
				}
			}
			loworderbits = loworderbits >> (32 - dictsize)
			offset =  (offset<<dictsize) + int(loworderbits)
			oldLen := len(outBytes)
			offset = oldLen - offset - 1
			newSlice := make([]byte, length)
			outBytes = append(outBytes, newSlice...)
			for i:=offset; i<offset+length; i++ {
				outBytes[oldLen+i] = outBytes[offset+i]
			}

		} else {
			// literal byte in next 8 bits, lsb first
			var byt byte
			for i:=0; i<8; i++ {
				bit = bytes[byteOff] & (1 << bitOff) != 0
				byt = byt >> 1
				if bit {
					byt |= 0x80
				}
				bitOff++
				if bitOff == 8 {
					bitOff = 0
					byteOff++
				}
			}
			outBytes = append(outBytes, byt)
		}
		if length == lengthEndOfStream {
			break
		}
	}
	return outBytes, nil
}

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
			return uncData.Bytes(), FileError{err}
		}
		switch tokenFlag {
		// bit 1 indicates length/offset sequences follow
		case true:
			length, err = civ3Bitstream.lengthsequence()
			if err != nil {
				return uncData.Bytes(), err
			}
			// log.Printf("length %v", length)
			// The token equating to length 519 is the end-of-stream token
			if length != lengthEndOfStream {
				// If length is 2, then only two low-order bits are read for offset
				dictsize := header[1]
				if length == 2 {
					dictsize = 2
				}
				offset, err := civ3Bitstream.offsetsequence(int(dictsize))
				if err != nil {
					return uncData.Bytes(), err
				}
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
						return uncData.Bytes(), FileError{err}
					}
					uncData.WriteByte(byt)
				}
			}
		// bit 0 inticates next 8 bits are literal byte, lsb first
		case false:
			{
				literalByte, err := civ3Bitstream.ReadByte()
				if err != nil {
					return uncData.Bytes(), FileError{err}
				}
				uncData.Write([]byte{literalByte})
			}
		}
	}
	return uncData.Bytes(), error(nil)

}

func (b *BitReader) lengthsequence() (int, error) {
	var sequence bitKey
	for {
		bit, err := b.ReadBit()
		if err != nil {
			return 0, FileError{err}
		}
		sequence.key = sequence.key << 1
		sequence.keyBitLength++
		if bit {
			sequence.key = sequence.key | 0b1
		}
		if sequence.keyBitLength > 8 {
			return 0, DecodeError{"Token did not match length sequence"}
		}
		if _, keyPresent := lengthLookup[sequence]; keyPresent {
			break
		}
	}
	xxxes, err := b.ReadBits(uint(lengthLookup[sequence].extraBits))
	if err != nil {
		return 0, FileError{err}
	}
	return lengthLookup[sequence].value + int(xxxes), nil
}
func (b *BitReader) offsetsequence(dictsize int) (int, error) {
	var sequence bitKey
	for {
		bit, err := b.ReadBit()
		if err != nil {
			return 0, FileError{err}
		}
		sequence.key = sequence.key << 1
		sequence.keyBitLength++
		if bit {
			sequence.key = sequence.key | 0b1
		}
		if sequence.keyBitLength > 8 {
			return 0, DecodeError{"Token did not match offset sequence"}
		}
		if _, keyPresent := offsetLookup[sequence]; keyPresent {
			break
		}
	}
	loworderbits, err := b.ReadBits(uint(dictsize))
	if err != nil {
		return 0, FileError{err}
	}
	return offsetLookup[sequence]<<uint(dictsize) + int(loworderbits), nil
}

// ReadByte reads a single byte from the stream, regardless of alignment
func (b *BitReader) ReadByte() (byte, error) {
	// If I init inside the loop this is out of scope
	var byt byte

	// Shift in 8 bits, LSBit first
	for i := 0; i < 8; i++ {
		bit, err := b.ReadBit()
		if err != nil {
			return 0, FileError{err}
		}
		byt >>= 1
		if bit {
			byt |= 128
		}
	}

	return byt, nil
}

// ReadBits reads  nbits from the stream
func (b *BitReader) ReadBits(nbits uint) (uint, error) {
	// If I init inside the loop this is out of scope
	var value uint

	// Use power to assign bits lsb to msb
	for i := uint(0); i < nbits; i++ {
		bit, err := b.ReadBit()
		if err != nil {
			return 0, FileError{err}
		}
		if bit {
			value |= 1 << i
		}
	}

	return value, nil
}
