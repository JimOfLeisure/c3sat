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
	// get the first four bytes to determine file type
	header := readBytes(r, 4)
	// reset pointer to parse from beginning
	r.Seek(0, 0)
	switch string(header) {
	case "CIV3":
		// log.Println("Civ3 save file detected")
		readcivheader(r)
		readbic(r)
	case "BIC ", "BICX":
		// log.Fatal("Civ3 BIC file detected. Currently not parsing these directly.")
		readbic(r)
	default:
		log.Fatalf("Civ3 file not detected. First four bytes:\n%s", hex.Dump(header))
	}
}

func check(e error) {
	if e != nil {
		// panic(e)
		log.Fatalln(e)
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
	name := readBytes(r, 4)
	length := readBytes(r, 4)
	c.name = string(name[:4])
	c.length = binary.LittleEndian.Uint32(length[:4])
	c.buffer.Write(readBytes(r, int(c.length)))
	return
}

func somethingsdifferent(s string, r *bytes.Reader) {
	// seeking backwards is causing EOF when I pass
	// r.Seek(-4, 2)
	log.Fatalf("%s\n%s\n", s, hex.Dump(readBytes(r, 1024)))
}

func readcivheader(r *bytes.Reader) {
	civ3header := readBytes(r, 30)
	_ = civ3header
	/*
	 * The first 30 bytes are identical on 151 saves on my PC, presumably all from c3c but some originally downloaded from CivFan
	 * But I downloaded 2 GOTMs, 2 COTMs and LK154 SG (running CCM mod), and they all have different info from mine and from each other
	 */
	// civ3headerref := [30]byte{0x43, 0x49, 0x56, 0x33, 0x00, 0x1a, 0x18, 0x00, 0x00, 0x00, 0x0a, 0x00, 0x00, 0x00, 0x7a, 0x1d, 0x96, 0xd6, 0x27, 0xca, 0x54, 0x4b, 0xa2, 0x76, 0x96, 0xd0, 0x81, 0x5d, 0x1a, 0xf7}
	// for i := 0; i < len(civ3headerref); i++ {
	// 	if civ3header[i] != civ3headerref[i] {
	// 		log.Fatal(hex.Dump(civ3header))
	// 	}
	// }
	// log.Println(hex.Dump(civ3header))
	bicheader := readBase(r)
	_ = bicheader
	// log.Println(bicheader.name, bicheader.length)
	// log.Println(hex.Dump(bicheader.buffer.Bytes()))

}

func readbic(r *bytes.Reader) {

	// log.Println(hex.Dump(readBytes(r, 1024)))

	// The next three values always seem to be the same in all files
	// Best guess is that the 1 is a record count and 720 is a length,
	// but almost all zeroes in the 720 bytes. But it does seem to lead
	// up to the next class name "GAME" (or "BLDG" in scenario games)
	bicqvernum := readBytes(r, 8)
	switch string(bicqvernum) {
	case "BICQVER#", "BICXVER#":
	default:
		log.Fatal(string(bicqvernum))
		// log.Fatal(somethingsdifferent(r))
		somethingsdifferent(string(bicqvernum), r)
	}
	bicone := readBytes(r, 4)
	if binary.LittleEndian.Uint32(bicone) != 1 {
		somethingsdifferent(string(binary.LittleEndian.Uint32(bicone)), r)
	}
	homelesslength := int(binary.LittleEndian.Uint32(readBytes(r, 4)))
	if homelesslength != 720 {
		somethingsdifferent(string(binary.LittleEndian.Uint32(bicone)), r)
	}
	homelessdata := readBytes(r, homelesslength)
	_ = homelessdata

	bicgame := readBytes(r, 4)
	if string(bicgame) != "GAME" {
		somethingsdifferent(string(bicgame), r)
	}
	// bicgame := readBase(r)
	// log.Println(bicgame.name, hex.Dump(bicgame.buffer.Bytes()))

	// Actually this does look like a list of integers, but it's more than 28

	// Per Antal1987's Bic data structure, it looks like 28 ints then some other data is at the start of the BIC
	// But the data doesn't seem to look like that
	// twentyeightints := readBytes(r, 28*4)
	// log.Println(hex.Dump(twentyeightints))

	// log.Println(hex.Dump(readBytes(r, 1024)))

}
