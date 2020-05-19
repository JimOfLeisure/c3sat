# civ3decompress

Civ3Decompress package decompresses Civ III save and bic/bix/biq files.

It is implemented based on the description of PKWare Data Compression Library at
https://groups.google.com/forum/#!msg/comp.compression/M5P064or93o/W1ca1-ad6kgJ
. However this is only a partial implementation; The Huffman-coded literals of
header 0x01 are not implemented here as they are not needed for decompressing
Civ3 data files.

It has worked fine on every Civ3 file I've used it on for years, and for
one-at-a-time decompressions it seems fast enough. But since I started Lua
scripting and batch processing I notice it's a relatively slow decompressor.

## Exports

- `func ReadFile(path string) ([]byte, bool, error)` - Most likely what you
want. Given a path to a Civ3 data file, it will detect whether or not it's
compressed and then return a byte array of the decompressed file contents.
- `func Decompress(file io.Reader) ([]byte, error)` - Given an `io.Reader` for
a compressed Civ3 data file, returns a byte array of the decompressed file.
- `type FileError struct` - If you want to handle errors based on type
- `type DecodeError struct` - If you want to handle errors based on type
- `func (b *BitReader) ReadByte() (byte, error)` - I don't
recall why this would be exported
- bitstream wrapper - These functions behave exactly like thier
[bitstream](https://github.com/dgryski/go-bitstream/blob/master/bitstream.go)
counterparts except they pop the least significant bit first. I'm not sure it
makes sense to export them.
  - `type BitReader struct`
  - `func NewReader(r io.Reader) *BitReader`
  - `func (b *BitReader) ReadBit() (Bit, error)`
  - `type Bit bool`
  - `type BitReader struct`
  - `const Zero Bit = false`
  - `const One = true`
