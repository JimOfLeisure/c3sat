package luaciv3

import (
	"io/ioutil"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

// CivString Finds null-terminated string and converts from Windows-1252 to UTF-8
func CivString(b []byte) (string, error) {
	var win1252 string
	var i int
	for i = 0; i < len(b); i++ {
		if b[i] == 0 {
			win1252 = string(b[:i])
			break
		}
	}
	if i == len(b) {
		win1252 = string(b)
	}
	sr := strings.NewReader(win1252)
	tr := charmap.Windows1252.NewDecoder().Reader(sr)

	outUtf8, err := ioutil.ReadAll(tr)
	if err != nil {
		return "", err
	}

	return string(outUtf8), nil
}

// Under development; pass it an offset for a list and a callback function
// It will call the callback with an offset for each element of the list
func listSection(offset int, callback func(off int)) {
	var itemLen int
	numItem := currentBic.readInt32(offset, Signed)
	// skip over count
	off := offset + 4
	for i := 0; i < numItem; i++ {
		itemLen = currentBic.readInt32(off, Signed)
		// skip over the length
		off += 4
		callback(off)
		off += itemLen
	}
}
