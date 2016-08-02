/*
This file's contents taken from https://github.com/dgryski/go-bitstream/blob/master/bitstream.go
and truncated for read-only and adapted to read the least-significant bit first which is backwards from the original code

Curiously, ReadByte didn't need to be reversed, just ReadBit (ReadBits remains to be seen)

License noticed copied from source. Only minor bit rotation changes made by me, Jim Nelson:

The MIT License (MIT)

Copyright (c) 2015 Damian Gryski <damian@gryski.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/

// Package bitstream is a simple wrapper around a io.Reader and io.Writer to provide bit-level access to the stream.
package main

import (
	"io"
)

// A Bit is a zero or a one
type Bit bool

const (
	// Zero is our exported type for '0' bits
	Zero Bit = false
	// One is our exported type for '1' bits
	One = true
)

// A BitReader reads bits from an io.Reader
type BitReader struct {
	r     io.Reader
	b     [1]byte
	count uint8
}

// NewReader returns a BitReader that returns a single bit at a time from 'r'
func NewReader(r io.Reader) *BitReader {
	b := new(BitReader)
	b.r = r
	return b
}

// ReadBit returns the next bit from the stream, reading a new byte from the underlying reader if required.
func (b *BitReader) ReadBit() (Bit, error) {
	if b.count == 0 {
		if n, err := b.r.Read(b.b[:]); n != 1 || (err != nil && err != io.EOF) {
			return Zero, err
		}
		b.count = 8
	}
	b.count--
	d := (b.b[0] & 0x01)
	b.b[0] >>= 1
	return d != 0, nil
}

// ReadByte reads a single byte from the stream, regardless of alignment
func (b *BitReader) ReadByte() (byte, error) {

	if b.count == 0 {
		n, err := b.r.Read(b.b[:])
		if n == 0 {
			b.b[0] = 0
		}
		return b.b[0], err
	}

	byt := b.b[0]

	var n int
	var err error
	n, err = b.r.Read(b.b[:])
	if n != 1 || (err != nil && err != io.EOF) {
		return 0, err
	}

	byt |= b.b[0] >> b.count

	b.b[0] <<= (8 - b.count)

	return byt, err
}

// ReadBits reads  nbits from the stream
func (b *BitReader) ReadBits(nbits int) (uint64, error) {

	var u uint64

	for nbits >= 8 {
		byt, err := b.ReadByte()
		if err != nil {
			return 0, err
		}

		u = (u << 8) | uint64(byt)
		nbits -= 8
	}

	var err error
	for nbits > 0 && err != io.EOF {
		byt, err := b.ReadBit()
		if err != nil {
			return 0, err
		}
		u <<= 1
		if byt {
			u |= 1
		}
		nbits--
	}

	return u, nil
}
