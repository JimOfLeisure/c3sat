Civ3 Show-And-Tell

This is something I started hacking on the past weekend, and it's to the
point I want to refactor it and organize it better. It is also nearing
the point it might be of interest to someone else, so I am putting it
on Github. I will be adding copyright notices and probably release this
under the GPL. At the moment I am trying to ensure I have everything
collected and uploaded and will try to keep it working as such while I
reorganize it.

This started as a collection of vague ideas, and I started posting in a
forum about it here: http://forums.civfanatics.com/showthread.php?t=493582

The idea and ultimate goal is for web browsers to be able to view
non-game-spoiling information from the Civilization III Conquests game
(such as might be seen in a screenshot) given a save game file (such
as might be posted to a forum for advice or succession game handoff)
without needing to open the game save in the game.

The short-term goal is to parse the save game file, extract interesting
information and display it in HTML format.

unc-test.sav is a decompressed game save file I am using during
development. When I run into trouble so far I have hard-coded to work
with this particular file.

civ3parse.py was my second effort at simply reading 4-char section
headers from the save file. It does not assume a data order and kept
getting tripped up by perceived inconsitencies in the game file.

tileonly.py is the most current and successful file I've been working
on. Its purpose is to simply read in the map tile data and optionally
generate an HTML map. It is hard-coded to read data from my test save. I
have been repeatedly altering it and running it directly as I attempt
to understand the game save data.

html.py calls tileonly.py and writes an html map in ./html/www/map.html .

To Do (short term):
- Copy map.css and map.js into repo; symlink from my dev web server to repo
- Figure out what to do about jquery.js; possibly link to Google's hosted jquery.js? Currently I have a copy on my local server
- Reorganize code so hack code is largely in separate module from the classes
- Implement debug and spoiler triggers to turn debug prints and game-spoiler info on/off
- Make tileonly export tile data in JSON and then have jquery build the HTML map table
- Break out my scant notes on save data into doc folder and organize it
- Include links to any reference info I've found on save file format

Medium-term goals:
- Get it to read any C3C save
- Figure a way to auto-decompress saves. blast()? Currently am using dynamite program to manually decompress save files.
- Enable with jQuery html map view where I can dynamically select which offset raw data is displayed in tile (to help visually figure out what each value means).
