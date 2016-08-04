package parseciv3

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"log"
)

type baseClass struct {
	name   string
	length uint32
	buffer bytes.Buffer
}

// Parseciv3 ...
func Parseciv3(civdata []byte) {
	r := bytes.NewReader(civdata)
	header := readBytes(r, 4)
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
	civ3header := readBytes(r, 30)
	log.Println(hex.Dump(civ3header))
	log.Printf("%v", readBase(r))
	log.Println(hex.Dump(readBytes(r, 16)))
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// readBytes repeatedly calls bytes.Reader.ReadByte()
func readBytes(r *bytes.Reader, n int) []byte {
	var out bytes.Buffer
	for i := 0; i < n; i++ {
		byt, err := r.ReadByte()
		check(err)
		out.WriteByte(byt)
	}
	return out.Bytes()
}

func readBase(r *bytes.Reader) (c baseClass) {
	// err := binary.Read(r, binary.LittleEndian, c.length)
	// check(err)
	name := readBytes(r, 4)
	length := readBytes(r, 4)
	// buffer := readBytes(r, length)

	// c = baseClass{
	// name:   zname,
	// 	length: zlength,
	// }
	c.name = string(name[:4])
	c.length = binary.LittleEndian.Uint32(length[:4])
	c.buffer.Write(readBytes(r, int(c.length)))
	return
}
