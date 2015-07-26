#!/usr/bin/env python
# -*- coding: latin-1 -*-

# Copyright (c) 2013, 2014 Jim Nelson
# 
# Permission is hereby granted, free of charge, to any person obtaining
# a copy of this software and associated documentation files (the
# "Software"), to deal in the Software without restriction, including
# without limitation the rights to use, copy, modify, merge, publish,
# distribute, sublicense, and/or sell copies of the Software, and to
# permit persons to whom the Software is furnished to do so, subject to
# the following conditions:
# 
# The above copyright notice and this permission notice shall be
# included in all copies or substantial portions of the Software.
# 
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
# EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
# MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
# NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
# LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
# OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
# WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

# 2014-05-19 Now I want to skip to the first WRLD section in arbitrary uncompressed SAV files

import struct	# For parsing binary data
import json     # to export JSON for the HTML browser
from horspool import horspool    # to seek to first match; from http://inspirated.com/2010/06/19/using-boyer-moore-horspool-algorithm-on-file-streams-in-python
import sys

# for debug; I am using this to change which civ visible_to I'm viewing
myglobalciv = 1
#myglobalcompare = 0xffff
myglobalcompare = 0
#myglobalbitmask = 0xefc0
#myglobalbitmask = 0xff2c
myglobalbitmask = 0x00ff

class GenericSection:
    """Base class for reading SAV sections."""
    def __init__(self, saveStream):
        buffer = saveStream.read(8)
        (self.name, self.length,) = struct.unpack_from('4si', buffer)
        #self.offset = saveStream.tell()
        self.buffer = saveStream.read(self.length)

class Tile:
    """Class for each logical tile."""
    def __init__(self, saveStream, debug=False):

        self.Tile36 = GenericSection(saveStream)
        self.rivers = get_byte(self.Tile36.buffer, 0x00)
        self.Tile36.values = struct.unpack_from('2h3i2c9h', self.Tile36.buffer)
        self.continent = self.Tile36.values[11]
        self.continent = self.Tile36.values[11]
        self.top_unit_id = self.Tile36.values[3]
        self.resource = self.Tile36.values[2]
        self.barb_info = self.Tile36.values[8]
        self.city_id = self.Tile36.values[9]
        self.colony = self.Tile36.values[10]
        #self.whatsthis = self.Tile4.values[0]
        self.whatsthis = get_short(self.Tile36.buffer, 0x00) & myglobalbitmask
        if not debug:
            del self.Tile36

        self.Tile12 = GenericSection(saveStream)
        self.Tile12.values = struct.unpack_from('iihh', self.Tile12.buffer)
        self.terrain = get_byte(self.Tile12.buffer, 0x5)
        self.improvements = self.Tile12.values[0]
        self.terrain_features = get_short(self.Tile12.buffer, 0x0a)
        # Mask 0x0001 is *** Bonus Grassland *** . Interesting, some hills and mountains have it (base tile is grassland)
        # Mask 0x0002 is the "fat x" around each city whether or not it has the culture to work it
        # Mask 0x0004 - Can't find any
        # Mask 0x0008 - Player start location
        # Mask 0x0010 - Snow-caps for mountains
        # Mask 0x0020 - Unsure. Only on land, seems clumped. Seems to be all tundra on generated maps and all forest tundra on LK's WM. Have not seen it on jungle tiles.
        # Other nybble masks 0x00c0 - nothing I can find
        # Mask 0x1000 - This looks like a likely candidate for "forest already chopped here"
        # Other nybble masks 0xe000 - nothing I can find ***  (found some of these on LK's WM)
        # Mask 0x0100 - river N corner? or NW?
        # Mask 0x0200 - river W corner?
        # Mask 0x0400 - river E corner? or SE?
        # Mask 0x0800 - river S corner?
        if not debug:
            del self.Tile12

        self.Tile4 = GenericSection(saveStream)
        self.Tile4.values = struct.unpack_from('i', self.Tile4.buffer)
        if not debug:
            del self.Tile4

        self.Tile128 = GenericSection(saveStream)
        self.Tile128.values = struct.unpack_from('4i96b', self.Tile128.buffer)

        self.is_visible_to_flags = self.Tile128.values[0]

        self.is_visible_now_to_flags = self.Tile128.values[1]

        self.worked_by_city_id = get_short(self.Tile128.buffer, 0x14)

        mytemp = 0x16
        self.land_trade_network_id = []
        for civ in range(32):
            self.land_trade_network_id.append(get_short(self.Tile128.buffer, mytemp))
            mytemp += 2

        mytemp = 0x56
        self.improvements_known_to_civ = []
        for civ in range(32):
            self.improvements_known_to_civ.append(get_byte(self.Tile128.buffer, mytemp))
            mytemp += 1

        if not debug:
            del self.Tile128

        # Remove values that are -1
        for key in self.__dict__.keys():
            if self.__dict__[key] == -1:
                del self.__dict__[key]

class Tiles:
    """Class to read all tiles"""
    def __init__(self, saveStream, width, height, debug=False):
        self.width = width      # These may eventually be redundant to a parent class
        self.height = height
        self.tile = []          # List of individual tiles
        #self.tile_matrix = []       # x,y matrix of individual tiles
        #self.tile_iso_matrix = []   # faux isometric padded matrix  of individual tiles
#        logical_tiles = width / 2 * height
#        while logical_tiles > 0:
#            self.tile.append(Tile(saveStream))
#            logical_tiles -= 1
        for y in range(height):
            #self.tile_matrix.append([])
            #self.tile_iso_matrix.append([])
            #if y % 2 == 1:
            #    self.tile_iso_matrix[y].append(None)
            for x in range(width / 2):
                this_tile = Tile(saveStream, debug)

                self.tile.append(this_tile)

                #self.tile_matrix[y].append(this_tile)

                #self.tile_iso_matrix[y].append(this_tile)
                #self.tile_iso_matrix[y].append(None)

            #if y % 2 == 0:
            #    self.tile_iso_matrix[y].append(None)

    def jsonDefault(self, o):
        """Trying to make json.dumps() work on all my data"""
        return o.__dict__

    def json_out(self, spoiler=False, debug=False):
        """Return a string of json-coded map"""
        #return json.dumps(self, default=self.jsonDefault, indent=4)
        return json.dumps(self, default=self.jsonDefault)

class Wrld:
    """Class for 3 WRLD sections"""
    def __init__(self, saveStream, debug=False):
        """Currently calling this from the horspool seek, so WRLD is already consumed from the stream. Read the length first."""
        self.name = "WRLD"
        buffer = saveStream.read(4)
        (self.length,) = struct.unpack_from('i', buffer)
        #print self.length
        self.buffer = saveStream.read(self.length)
        # Extract any data here, but I think it's only 2 bytes
        (self.num_continents,) = struct.unpack_from('h', self.buffer)
        #print self.name
        #print hexdump(self.buffer)
        if not debug:
            del self.buffer

        self.Wrld2 = GenericSection(saveStream)
        self.Wrld2.values = struct.unpack_from('41i', self.Wrld2.buffer)
        self.height = self.Wrld2.values[1]
        self.width = self.Wrld2.values[6]
        # Civ start locations in the form of map tile index numbers
        self.start_loc = []
        for myindex in range(7,39):
            self.start_loc.append(self.Wrld2.values[myindex])
        self.world_seed = self.Wrld2.values[39]
        self.ocean_continent_id = self.Wrld2.values[0]
            
        if not debug:
            del self.Wrld2

        self.Wrld3 = GenericSection(saveStream)
        self.Wrld3.values = struct.unpack_from('13i', self.Wrld3.buffer)

        if not debug:
            del self.Wrld3

        self.Tiles = Tiles(saveStream, self.width, self.height, debug)

        # Total hack; I want this info in svg_out, so I'm stuffing it into Tiles rather than refactor
        self.Tiles.start_loc = self.start_loc

        # CONT sections
        self.continents = []
        for i in range(self.num_continents):
            # First integer is 0 if water, 1 if land
            # Second integer is number of tiles on the continent
            # index is the continent ID (0..num_continents-1)
            self.continents.append(struct.unpack_from('ii',(GenericSection(saveStream).buffer)))

def get_byte(buffer, offset):
    """Unpack an byte from a buffer at the given offest."""
    (the_byte,) = struct.unpack('B', buffer[offset:offset+1])
    return the_byte

def get_short(buffer, offset):
    """Unpack an int from a buffer at the given offest."""
    (the_short,) = struct.unpack('H', buffer[offset:offset+2])
    return the_short

def get_int(buffer, offset):
    """Unpack an int from a buffer at the given offest."""
    (the_int,) = struct.unpack('I', buffer[offset:offset+4])
    return the_int

def parse_save(saveFile, debug=False):
    buffer = saveFile.read(4)
    if buffer <> 'CIV3':
        print "wah wah wah wahhhhhhhh."
        print "Stub. Provided stream not decompressed C3C save"
        return -1
    #print 'Using Horspool search to go to first WRLD section'
    wrldOffset = horspool.boyermoore_horspool(saveFile, "WRLD")
    #print wrldOffset
    game = Wrld(saveFile, debug)
    saveFile.close()
    return game

def hexdump(src, length=16):
    """Totally yoinked from https://gist.github.com/sbz/1080258"""
    FILTER = ''.join([(len(repr(chr(x))) == 3) and chr(x) or '.' for x in range(256)])
    lines = []
    for c in xrange(0, len(src), length):
        chars = src[c:c+length]
        hex = ' '.join(["%02x" % ord(x) for x in chars])
        printable = ''.join(["%s" % ((ord(x) <= 127 and FILTER[ord(x)]) or '.') for x in chars])
        #lines.append("%04x  %-*s  %s\n" % (c, length*3, hex, printable))
        lines.append("%04x  %-*s  %s" % (c, length*3, hex, printable))
    return '\n'.join(lines)

def main():
    saveFile = open("gamesaves/unc-test.sav", 'rb')
    #saveFile = open("gamesaves/unc-lk151-650ad.sav", 'rb')
    game = parse_save(saveFile)

    print 'Printing something(s) from the class to ensure I have what I intended'
    #print game.name, game.length
    #print game.tile.pop()[1].length
    print "Map width:",game.width,"Map height:", game.height
    print game.Tiles.tile[0].info['terrain']
    print game.Tiles.tile[1000].info['terrain']
    #print game.tile[0].is_visible_to
    #max = len(game.tile)
    max = 10
    for i in range(max):
        print game.Tiles.tile[i].continent
    #game.Tiles.test_things()
    print game.start_loc
        

if __name__=="__main__":
    main()
