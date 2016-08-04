## Civ3 Show-And-Tell

August 2016 update: I'm embarking on another major rewrite/refactor. I have an interest in learning Go,
and this project would benefit greatly from using Go instead of Python. I could have better modularized
and managed the Python version, but the two big problems were that target users aren't likely to have
Python installed, and I was unable to reasonably compile the code into an executable. Go will let me
release native executables for Windows, Mac and Linux.

So far, all the new activity is in readciv3/, and instead of immediately recoding existing functionality
I'm writing a decompressor so external utilities won't be needed.

And now that the decompressor code is working I'll be refactoring into the civ3sat/ folder.

### To Do

- ~~Finish decompressor code~~
- ~~Refactor decompressor code into package~~ working in decompressor, but not pretty. Presume will eventually pass file/stream pointers around.
- ~~Use cli package to create command line app~~ Working in civ3sat/.
- Start extracting data from file - about to begin in parseciv3/.
- Clean up error checking an handling.
    - Get rid of log.Fatal
    - Bubble errors up
    - Handle in main
    - Figure out verbose/debug/quiet switches
- Flesh out subcommands like ~~decompress, hexdump,~~ info (including world random seed and game/map settings)
- Extract maps like Python code

## Tagged older versions

Both are in Python with an HTML/JavaScript viewer page. Lost to history is a C# version of an early save game parser, but it didn't have useful output.

- [0.1](https://github.com/myjimnelson/c3sat/tree/0.1) - This version output the SVG map to file and had a jQuery/HTML viewer page to view it.
- [0.2](https://github.com/myjimnelson/c3sat/tree/0.2) - This version output JSON map data, and a d3.js/HTML viewer rendered the map based on the JSON.

-----
## Readme from July 2015

This is an attempt at an accessory map viewer for Civilization III Conquests save game files.

Run `civ3tojson.py` and pass it a filename of an uncompressed save (use an autosave or decompress with dynamite) or pipe it an uncompressed save file (autosave or blast) and it will print JSON output which can be redirected to html/civmap.json and viewed as a map with html/d3.html.

Status: Currently changing from generating the SVG from Python to generating JSON output and constructing the map with d3.js in-browser. [Tag 0.1](https://github.com/myjimnelson/c3sat/tree/0.1) was the most public-ready version of the svg-out code.

I am licensing my "artwork" which includes SVG representations of mountains, hills and trees under a [Creative Commons Attribution 4.0 International License](http://creativecommons.org/licenses/by/4.0/). Don't try too hard; in your attribution you can link back to my GitHub repo, my GitHub user page or my CivFantatics Forums user page or a thread started by me; whatever is easy. Or heck just use my name. I'm just licensing it to assure you you can use it.

Some code from other authors is currently included in this repo. See [the horspool.py readme](horspool/readme.md) for license and attribution.

## History

This is something I started hacking on [in April 2013], and it's to the
point I want to refactor it and organize it better. It is also nearing
the point it might be of interest to someone else, so I put it
on Github.

This started as a collection of vague ideas, and I started posting in a
forum about it here: http://forums.civfanatics.com/showthread.php?t=493582

The idea and ultimate goal is for web browsers to be able to view
non-game-spoiling information from the Civilization III Conquests game
(such as might be seen in a screenshot) given a save game file (such
as might be posted to a forum for advice or succession game handoff)
without needing to open the game save in the game.

The short-term goal is to parse the save game file, extract interesting
information and display it in HTML/SVG format. The base and overlay
terrain is represented by SVG shapes or temporary text.

civ3parse.py is the most current and successful file I've been working
on. Its purpose is to read in the game data and output it in JSON format.
It currently uses horspool.py to find the first instance of WRLD in the
save game file and then proceed to read the map size and map tiles.

My original plan was to use JSON as an intermidate data format between the
save game file and the map display. At first it was easier to generate
the SVG from Python, but now that I'm learning d3.js I am returning to
the original idea.

## To Do (short term)
- Implement debug and spoiler triggers to turn debug prints and game-spoiler info on/off - Successfully omitting undiscovered tiles from the map. Currently do not have the data to understand which bonus/luxury/strategic resource locations are spoiler, so currently only identifying those which appear at the beginning of a default epic game.
- Include links to any reference info I've found on save file format
- Get it to read any C3C save - Mostly done. It's seeking and matching the 4-character labels in the file which apparently mark the start of C++ classes. But it doesn't natively decompress, and it doesn't understand the whole file.
- Figure a way to auto-decompress saves - Currently using blast program compiled from zlib contrib to manually decompress save files to stdout. Previously was using ubuntu package dynamite to decompress the file.

## Medium-term goals

## Long-term goals
- Python app server that will pull a save game file by URL, parse it and provide non-game-spoiling map and info as SVG to modern web browser. - Sort of done, but not real cohesive.
- Possibly allow map annotation (city-planning "dot maps", etc.) via JSON files generated by HTML page and posted to e.g. forums (Think this idea is a non-starter, but I'll leave it here for now...maybe use url queries to place user data?)
