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
	log.Fatalf("%s\n%s\n", s, hex.Dump(readBytes(r, 256)))
}

func readcivheader(r *bytes.Reader) {
	civ3header := readBytes(r, 30)
	_ = civ3header

	// Seems to be four or five ints followed by 2 strings, relative paths to bic resources and bic file
	bicheader := readBase(r)
	_ = bicheader
}

func readbic(r *bytes.Reader) {
	bicqvernum := readBytes(r, 8)
	switch string(bicqvernum) {
	case "BICQVER#", "BIC VER#", "BICXVER#":
	default:
		log.Fatal(string(bicqvernum))
		somethingsdifferent(string(bicqvernum), r)
	}
	// always seems to be 1
	bicone := readBytes(r, 4)
	if binary.LittleEndian.Uint32(bicone) != 1 {
		somethingsdifferent(string(binary.LittleEndian.Uint32(bicone)), r)
	}
	bicdescriptionlength := int(binary.LittleEndian.Uint32(readBytes(r, 4)))
	if bicdescriptionlength != 720 {
		somethingsdifferent(string(binary.LittleEndian.Uint32(bicone)), r)
	}
	bicdescription := readBytes(r, bicdescriptionlength)
	_ = bicdescription
	// log.Println(hex.Dump(bicdescription))
	// log.Println(hex.Dump(bicdescription[:16]))
	// log.Println(hex.Dump(bicdescription[8:12]))

	// At this point, my epic SAVs have GAME, downloaded scenario-based games have BLDG

	var buffernext []byte
	bicnext := string(readBytes(r, 4))
	switch bicnext {
	case "BLDG":
		log.Println(bicnext)
		numbuildings := int(binary.LittleEndian.Uint32(readBytes(r, 4)))
		for i := 0; i < numbuildings; i++ {
			buffernext = readBytes(r, 0x110)
			_ = buffernext
			// log.Println(hex.Dump(buffernext))
		}
		// CTZN
		bicnext = string(readBytes(r, 4))
		log.Println(bicnext)
		numcitizentypes := int(binary.LittleEndian.Uint32(readBytes(r, 4)))
		for i := 0; i < numcitizentypes; i++ {
			buffernext := readBytes(r, 0x80)
			_ = buffernext
		}
		// CULT
		bicnext = string(readBytes(r, 4))
		log.Println(bicnext)
		numcult := int(binary.LittleEndian.Uint32(readBytes(r, 4)))
		_ = numcult
		for i := 0; i < numcult; i++ {
			buffernext := readBytes(r, 0x5c)
			_ = buffernext
		}
		// DIFF
		bicnext = string(readBytes(r, 4))
		log.Println(bicnext)
		numdiff := int(binary.LittleEndian.Uint32(readBytes(r, 4)))
		_ = numdiff
		for i := 0; i < numdiff; i++ {
			// for i := 0; i < 1; i++ {
			buffernext := readBytes(r, 0x7c)
			_ = buffernext
		}
		// ERAS
		bicnext = string(readBytes(r, 4))
		log.Println(bicnext)
		numeras := int(binary.LittleEndian.Uint32(readBytes(r, 4)))
		_ = numeras
		for i := 0; i < numeras; i++ {
			buffernext := readBytes(r, 0x10c)
			_ = buffernext
		}
		// ESPN
		bicnext = string(readBytes(r, 4))
		log.Println(bicnext)
		numespn := int(binary.LittleEndian.Uint32(readBytes(r, 4)))
		_ = numespn
		for i := 0; i < numespn; i++ {
			buffernext := readBytes(r, 0xec)
			_ = buffernext
		}
		// EXPR
		bicnext = string(readBytes(r, 4))
		log.Println(bicnext)
		numexpr := int(binary.LittleEndian.Uint32(readBytes(r, 4)))
		_ = numexpr
		for i := 0; i < numexpr; i++ {
			// for i := 0; i < 1; i++ {
			buffernext := readBytes(r, 0x2c)
			_ = buffernext
			// log.Println(hex.Dump(buffernext))
		}
	default:
		log.Println("Unexpected class name: ", bicnext)

	}

	bicnext = string(readBytes(r, 4))
	log.Println(bicnext)
	log.Println(binary.LittleEndian.Uint32(readBytes(r, 4)))

	// var bicgame baseClass
	// bicgame.length = 1
	// for bicgame.name != "GAME" && 0 < bicgame.length && bicgame.length < 500 {
	// 	bicgame = readBase(r)
	// 	log.Println(bicgame.name, bicgame.length)
	// 	// log.Println(bicgame.name, hex.Dump(bicgame.buffer.Bytes()))
	// }

	log.Println(hex.Dump(readBytes(r, 0x200)))
	log.Println("")

}
