# This repo's develop branch

## Build notes

### CIA3

- Run "[pkger](https://github.com/markbates/pkger) -o /cmd/cia3" from within this repo before building to embed the html/* files into the executable, otherwise the exe will try to serve files from ./html/ and probably not find them.
- To get rid of the console window on Windows, build with `go build <path/to/cmd/cia3> -ldflags="-H windowsgui"`

## A "Map" of the repo

This repo now provides two executable programs. The executables either primarily or optionally start a GraphQL server on localhost for access to Civ3 save data.

- civ3sat in cmd/civ3sat
  - Civ 3 Show-and-Tell. The original intent was to generate sharable maps from Civ3 save files. In its latest versioned release it's a command-line utility to decompress or query Civ3 save files.
- cia3 in cmd/cia3
    - Civ Intelligence Agency III (CIA3) is intended to be a non-spoiling, non-cheating game assistant for single-human-player games of Sid Meier's Civilization III Complete and/or Conquests. Multiplayer non-spoiling support may be possible in the future.

The Go packages:

Somewhere in here there is/was some REST API server code. I don't see it in the recent branches, so I may have deleted it in favor of the GraphQL server code.
The GraphQL code is much more versatile.

- parseciv3 in parseciv3/* : Originally intended to read the Civ3 save file data into defined structures, but it's not all the way there.
- civ3decompress in civ3decompress/ implements an i/o reader for compressed Civ3 files. Its ReadFile will auto-detect whether or not a Civ3 file is compressed.
  - Decompress is implemented based on the description of PKWare Data Compression Library at https://groups.google.com/forum/#!msg/comp.compression/M5P064or93o/W1ca1-ad6kgJ
  - However this is only a partial implementation; The Huffman-coded literals of header 0x01 are not implemented here as they are not needed for my purpose
- queryciv3 in queryciv3/* : This is the main workhorse of the current versions of executables. It implements a GraphQL query api server to pull data from Civ3 saves by seeking to known offsets from named or known reference points.

Non-Go code:

- CIA3's UI is web-based, and its code is in cmd/cia3/html and is served at / on localhost by cia3.
  - This may include code from jsdiff and diff2html, at least in the develop branch
- Other HTML / JS code is in develop/html. It generated a map based on the save file, but it is not updated for the GraphQL API.
- Python 2: This repo originally started as Python, and the abandonedpydevelop branch points to the last dev version of that code. (Not in this develop branch)

Notes/documentation:

- develop/ has a collection of notes from myself and other places about the Civ3 save file format. It is not very organized.

## Future direction

At the moment, work is going into cmd/cia3 and queryciv3 with the intent of making a usable game assistant for others to use.

Taking a fresh look at this repo, here are some things that might be done, in no particular order:

- civ3sat
    - ~~Move civ3sat/ to cmd/civ3sat/~~ done
    - Update past map renderers to work with gql and have it make maps again
    - See about setting up a public upload/map server
    - Actually, not sure it needs to exist apart from cia3 as the work is being done by the GQL server and JavaScript
- ~~civ3satgql~~ queryciv3
    - ~~Rename civ3satgql to...something like civ3query, civ3gql, gqlciv3...to emphasize it's a query engine~~ done
    - Move servers out and just provide http handlers and support functions
    - ~~Add 'native' ability to read in file, wean off of parseciv3~~ just moved ReadFile() to civ3decompress which is a needed package, anyway
- parseciv3
    - ~~quit using it for auto-detecting compression and reading files when not using its other functions~~ done
    - leave it in place / in mothballs in case I want to take another crack at reading direct to structures later

## Wish list

Stuff I'd like to find and add.

- Civ attitude (polite/cautious/annoyed/etc.)
- Culture opinion both directions ("The Indian people are impressed with our culture")
- Best known enemy unit
- Trades
  - Techs (next priority)
  - Gold
  - Resources
  - Workers
  - Contacts (hide column until tech is known?)
  - Maps (hide column until tech is known?)
- Alerts
  - Civ at war will talk
  - New contact
  - New tech trades available
  - Unhappiness