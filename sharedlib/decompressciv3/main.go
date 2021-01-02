package main

import (
	"C"
	"fmt"
	"unsafe"

	"github.com/myjimnelson/c3sat/civ3decompress"
)
import "encoding/binary"

//export HelloDll
func HelloDll() {
	fmt.Println("Hello frm the decompressciv3 shared library.")
}

// Given a string path, returns the decompressed byte array if compressed or raw file othrwise
// Prepends an 8-byte little-endian byte length of the byte array payload
//export ReadFile
func ReadFile(path *C.char) unsafe.Pointer {
	outBytes, _, err := civ3decompress.ReadFile(C.GoString(path))
	if err != nil {
		panic(err)
	}
	length := make([]byte, 8)
	binary.LittleEndian.PutUint64(length, uint64(len(outBytes)))
	return C.CBytes(append(length, outBytes...))
}

// Given a compressed byte array and its length, returns the decompressed byte array
// Prepends an 8-byte little-endian byte length of the byte array payload
//export Decompress
func Decompress(data *C.char, length C.int) unsafe.Pointer {
	goData := C.GoBytes(unsafe.Pointer(data), length)
	outBytes, err := civ3decompress.DecompressByteArray(goData)
	if err != nil {
		panic(err)
	}
	outLength := make([]byte, 8)
	binary.LittleEndian.PutUint64(outLength, uint64(len(outBytes)))
	return C.CBytes(append(outLength, outBytes...))
}

func main() {
	fmt.Print("This code was meant to be compiled as a shared libary, not an executable.")
}
