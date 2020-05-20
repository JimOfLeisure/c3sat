package civ3decompress

import (
	"bytes"
	"crypto/sha1"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestDecompress(t *testing.T) {
	aiSha1Sum := []byte{174, 217, 185, 252, 134, 178, 180, 122, 76, 70, 5, 61, 245, 63, 235, 84, 78, 207, 49, 206}
	// get filename of current file; will use relative path from here for test data input
	_, filename, _, _ := runtime.Caller(0)
	path := filepath.Dir(filename) + "/_test_data/AI.bin"
	// Open file, hanlde errors, defer close
	file, err := os.Open(path)
	if err != nil {
		t.Fatal("Error opening test file:", err)
	}
	defer file.Close()
	data, err := Decompress(file)
	if err != nil {
		t.Fatal("Decompress:", err)
	}
	sum := sha1.Sum(data)
	if !bytes.Equal(sum[:], aiSha1Sum) {
		t.Errorf("Sha1 sum of output doesn't match. Expected %v, got %v", aiSha1Sum, sum)
	}
}

func TestDecompressByteArray(t *testing.T) {
	aiSha1Sum := []byte{174, 217, 185, 252, 134, 178, 180, 122, 76, 70, 5, 61, 245, 63, 235, 84, 78, 207, 49, 206}
	// get filename of current file; will use relative path from here for test data input
	_, filename, _, _ := runtime.Caller(0)
	path := filepath.Dir(filename) + "/_test_data/AI.bin"
	inData, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal("Error reading test file:", err)
	}
	outData, err := DecompressByteArray(inData)
	if err != nil {
		t.Fatal("Decompress:", err)
	}
	sum := sha1.Sum(outData)
	if !bytes.Equal(sum[:], aiSha1Sum) {
		t.Errorf("Sha1 sum of output doesn't match. Expected %v, got %v", aiSha1Sum, sum)
	}
}
