# parseciv3

This package was an attempt to read Civ III save files directly into Go structs.

None of my code uses this anymore, but I often refer to it when querying data.

I've decided it's much easier just to grab data by seeking to the 4-character
"section headers" and pulling offsets from that rather than try to properly
parse the file directly into data structures.