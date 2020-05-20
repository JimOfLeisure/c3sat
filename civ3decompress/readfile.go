package civ3decompress

import (
	"io/ioutil"
	"os"
)

// ReadFile takes a filename and returns the decompressed file data or the raw data if it's not compressed. Also returns true if compressed.
func ReadFile(path string) ([]byte, bool, error) {
	// Open file, hanlde errors, defer close
	file, err := os.Open(path)
	if err != nil {
		return nil, false, FileError{err}
	}
	defer file.Close()

	var compressed bool
	var data []byte
	header := make([]byte, 2)
	_, err = file.Read(header)
	if err != nil {
		return nil, false, FileError{err}
	}
	// reset pointer to parse from beginning
	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, false, FileError{err}
	}
	data, err = ioutil.ReadFile(path)
	if err != nil {
		return nil, false, FileError{err}
	}
	if header[0] == 0x00 && (header[1] == 0x04 || header[1] == 0x05 || header[1] == 0x06) {
		compressed = true
		uncData, err := DecompressByteArray(data)
		if err != nil {
			return nil, false, err
		}
		data = uncData
	}
	return data, compressed, error(nil)
}
