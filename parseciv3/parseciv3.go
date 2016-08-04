package parseciv3

import (
	"bytes"
	"encoding/hex"
	"log"
)

// Parseciv3 ...
func Parseciv3(civdata []byte) {
	r := bytes.NewReader(civdata)
	header := ReadBytes(r, 4)
	switch string(header) {
	case "CIV3":
		log.Println("Civ3 save file detected")
	case "BIC ", "BICX":
		// TODO: Intelligently parse BIC if the file is a BIC
		log.Fatal("Civ3 BIC file detected. Currently not parsing these directly.")
	default:
		log.Fatalf("Civ3 file not detected. First four bytes:\n%s", hex.Dump(header))
	}
	r.Seek(0, 0)
	civ3header := ReadBytes(r, 30)
	log.Println(hex.Dump(civ3header))
	log.Println(hex.Dump(ReadBytes(r, 4)))
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// ReadBytes repeatedly calls bytes.Reader.ReadByte()
func ReadBytes(r *bytes.Reader, n int) []byte {
	var out bytes.Buffer
	for i := 0; i < n; i++ {
		byt, err := r.ReadByte()
		check(err)
		out.WriteByte(byt)
	}
	return out.Bytes()
}
