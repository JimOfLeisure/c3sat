Civ3 Show-And-Tell

Tagged previous effort as 0.1 and created a develop branch. I now want
to store the game data in a database, so there will be code to read
the game file and insert data into a database, and then any spoiler,
map or report logic will happen after that. I also want a REST API for
accessing data from the DB.

Things are about to get broken in the develop branch.

I need to remind myself what each file does:
- civ3parse.py - One of the earlier attempts at reading the game file. It just reads and skips sections and prints debug info
- svg.py - The most current script. It reads an uncompressed save file passed as argument or as stdin, passes it through wrld.py and writes an svg file.
- tileonly.py - The first script to output map data in graphic or fancy text form. It's hard-coded to skip to the first TILE of my test save and parse the map info. It has several ouptut functions to write the map in various formats.
- webapp.py - My first attempt at making this easier to use by others. It is meant to be run from uWSGI, take a "url=" parameter which points to an on-the-internet save game file, fetch that save, decompress it with blast (from zlib contrib) if necessary and return an SVG of the map. It uses WRLD and it works.
- whatsthis.py - Uses wrld.py to parse map and prints or compares tile info. I kept needing to check on different data while searching for info in the save, and I kept editing this to do it.
- wrld.py - This is--up to now--the main program functions. Its parse_save function takes an uncompressed game file stream as input and loads it into Python classes. Its svg_out function takes that data and returns an SVG map. It uses a search function to find the first WRLD section and can pull the map width and height from it, so it works on arbitrary decompressed game save files.

Note to self: The copied file doc/civ3/Civ3save/Civ3sav/save.txt has some handy info I'll need soon
