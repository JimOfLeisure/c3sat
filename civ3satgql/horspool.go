package civ3satgql

// Adapting an MIT-licensed Horspool search implementation from https://raw.githubusercontent.com/callmekatootie/strsearch/master/strsearch.go
// I (Jim Nelson) will alter it for this module's needs
// Will search byte arrays for four capital ascii codes (plus another character or two) repeatedly

/*
The MIT License (MIT)

Copyright (c) 2015 Mithun Kamath

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

/*
   Public function that takes in the text and the pattern to find in that text
   Returns -1 if pattern can't be found, else returns the location (zero based)
   where pattern can be found
   Uses Boyer Moore Horspool algorithm
*/
func Find(text, pattern string) int {
	m := len(text)
	n := len(pattern)

	if n == 0 {
		return 0
	} else if m == n && text == pattern {
		return 0
	}

	//Prepare the table
	jumps := make(map[rune]int)

	for i, char := range pattern {
		if i < n-1 {
			jumps[char] = n - i - 1
		} else {
			// Mismatch of the last character results in a jump
			// equal to the length of the string
			jumps[char] = n
		}
	}

	j := n - 1

	for j < m {
		k := j
		i := n - 1

		for i >= 0 {
			if text[k] != pattern[i] {
				if jump, ok := jumps[rune(text[j])]; ok {
					// Character exists in the pattern
					// Jump based on the table created earlier
					j += jump
				} else {
					// Character does not exist in the pattern.
					// Jump the length of the pattern.
					j += n
				}

				break
			}

			k--
			i--
		}

		if i < 0 {
			// The pattern string has exhausted - We have a match
			return k + 1
		}
	}

	return -1
}
