# parseciv3

This package was an attempt to read Civ III save files directly into Go structs.

None of my code uses this anymore, but I often refer to it when querying data.

I've decided it's much easier just to grab data by seeking to the 4-character
"section headers" and pulling offsets from that rather than try to properly
parse the file directly into data structures.

It's been a few years since I've used it. I seem to recall it was incomplete at
best, but it would read up to a point, or maybe it didn't work on all saves.

## Exports

- `func NewCiv3Data(path string) (Civ3Data, error)` - Takes a path to a file and
returns a struct containing the parsed data and a rawFile field
- `func ParseCiv3(r io.ReadSeeker) (ParsedData, error)` - Takes raw save file
data and returns a map of the parsed data
- `type ParsedData map[string]Section`
- `type Civ3Data struct`
- Also categorized error types and lots of other data types