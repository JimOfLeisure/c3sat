#!/usr/bin/env python
# -*- coding: latin-1 -*-

# Copyright (c) 2013, 2014, 2015 Jim Nelson
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

import re	# Regular expressions
import struct	# For parsing binary data
import json     # to export JSON
from horspool import horspool    # to seek to first match; from http://inspirated.com/2010/06/19/using-boyer-moore-horspool-algorithm-on-file-streams-in-python
import sys

class Tiles:
    """Class to read all tiles"""
    def __init__(self, saveStream, width, height, debug=False):
        self.width = width      # These may eventually be redundant to a parent class
        self.height = height
        self.tile = []          # List of individual tiles
        self.tile_matrix = []       # x,y matrix of individual tiles
        self.tile_iso_matrix = []   # faux isometric padded matrix  of individual tiles

        # Tiles is pretty good, but time-consuming. Let's skip most of them while figuring out the rest of the file
        numTilesMinusTwo = (width / 2 * height) - 2
        # size of the tile sections for each tile in bytes
        myTileSize = 128 + 36 + 4 + 12 + (4*8)
        self.firstTile = Tile(saveStream, debug)
        # skip over most of the tiles
        saveStream.seek((numTilesMinusTwo * myTileSize), 1)
        self.lastTile = Tile(saveStream, debug)
        return

#        logical_tiles = width / 2 * height
#        while logical_tiles > 0:
#            self.tile.append(Tile(saveStream))
#            logical_tiles -= 1

        for y in range(height):
            self.tile_matrix.append([])
            self.tile_iso_matrix.append([])
            if y % 2 == 1:
                self.tile_iso_matrix[y].append(None)
            for x in range(width / 2):
                this_tile = Tile(saveStream, debug)

                self.tile.append(this_tile)

                self.tile_matrix[y].append(this_tile)

                self.tile_iso_matrix[y].append(this_tile)
                self.tile_iso_matrix[y].append(None)

            if y % 2 == 0:
                self.tile_iso_matrix[y].append(None)

class Wrld:
    """Class for 3 WRLD sections"""
    def __init__(self, saveStream, debug=False):
        self.Wrld1 = GenericSection(saveStream)
        # Extract any data here, but I think it's only 2 bytes
        (self.num_continents,) = struct.unpack_from('h', self.Wrld1.buffer)
        if not debug:
            del self.Wrld1.buffer

        self.Wrld2 = GenericSection(saveStream)
        #print self.Wrld2.name
        self.Wrld2.values = struct.unpack_from('41i', self.Wrld2.buffer)
        self.height = self.Wrld2.values[1]
        self.width = self.Wrld2.values[6]
        # Civ start locations in the form of map tile index numbers
        self.start_loc = []
        for myindex in range(7,39):
            self.start_loc.append(self.Wrld2.values[myindex])
        self.world_seed = self.Wrld2.values[39]
        self.ocean_continent_id = self.Wrld2.values[0]
            
        #print self.height
        #print self.width
        #print self.values
        #print hexdump(self.Wrld2.buffer)
        if not debug:
            del self.Wrld2.buffer
            del self.Wrld2.values

        self.Wrld3 = GenericSection(saveStream)
        self.Wrld3.values = struct.unpack_from('13i', self.Wrld3.buffer)
        #print self.Wrld3.name
        #print hexdump(self.Wrld3.buffer)
        if not debug:
            del self.Wrld3.buffer
            del self.Wrld3.values

        self.Tiles = Tiles(saveStream, self.width, self.height, debug)

        # Total hack; I want this info in svg_out, so I'm stuffing it into Tiles rather than refactor
        self.Tiles.start_loc = self.start_loc

        # CONT sections
        self.continents = []
        for i in range(self.num_continents):
            # First integer is 0 if water, 1 if land
            # Second integer is number of tiles on the continent
            # index is the continent ID (0..num_continents-1)
#        my_name = 'CONT'
#        my_count = 0
#        while my_name == 'CONT':
#            my_temp = GenericSection(saveStream)
#            my_name = my_temp.name
#            print my_count, my_temp.name, my_temp.length
#            print hexdump(my_temp.buffer)
#            my_count +=1
            self.continents.append(struct.unpack_from('ii',(GenericSection(saveStream).buffer)))

class Bic:
    """Class for the embedded BIC"""
    def __init__(self, saveStream, debug=False):
        self.offset = saveStream.tell()
        (self.name, self.verNumText, self.verNum, self.length,) = struct.unpack_from('4s4sii', saveStream.read(16))
        self.buffer = saveStream.read(self.length)
        self.hexdump = hexdump(self.buffer)

        #self.Game = GenericSection(saveStream)
        (self.GameName, self.GameVerNum, self.GameLength,) = struct.unpack_from('4sii', saveStream.read(12))
        self.GameHexdump = hexdump(saveStream.read(self.GameLength))

        #self.Game2 = GenericSection(saveStream)
        #(self.Game2Name, self.Game2VerNum, self.Game2Length,) = struct.unpack_from('4sii', saveStream.read(12))
        #self.Game2Hexdump = hexdump(saveStream.read(self.Game2Length))

        if not debug:
          del self.buffer
          #del self.whatsThis.buffer

class Civ3:
    """Class for entire save game"""
    def __init__(self, saveStream, debug=False):
        self.offset = saveStream.tell()
        self.name = saveStream.read(4)
        if self.name <> 'CIV3':
            print "wah wah wah wahhhhhhhh."
            print "Stub. Provided stream not decompressed C3C save"
            return -1
        self.buffer = saveStream.read(26)
        self.hexdump = hexdump(self.buffer, 26)
        # This Bic appears to be a save game class that contains a copy of a BIx file
        self.Bic = GenericSection(saveStream)
        self.Bic.Bic = Bic(saveStream, debug)

        ### Skipping 2nd GAME section as I can't figure out what it is. Think I'm also skipping a couple of DATEs and a PLGI

        horspoolOffset = horspool.boyermoore_horspool(saveStream, "CNSL")
        self.Cnsl = GenericSection(saveStream, "CNSL")
        self.Cnsl.horspoolOffset = horspoolOffset 

        self.Wrld = Wrld(saveStream, debug)

        ### Skipping some padding or other unknown data

        horspoolOffset = horspool.boyermoore_horspool(saveStream, "LEAD")
        self.Lead1 = GenericSection(saveStream, "LEAD")
        self.Lead1.horspoolOffset = horspoolOffset 

        #self.whatsNext = GenericSection(saveStream)
        self.whatsNextOffset = saveStream.tell()
        self.whatsNext = hexdump(saveStream.read(40000))
        #self.whatsNext = GenericSection(saveStream)

        if not debug:
          del self.buffer
          del self.Bic.buffer
          del self.Cnsl.buffer
          del self.Lead1.buffer
          #del self.whatsNext.buffer
        saveStream.close()

class newParse:
  """Starting over with parsing strategy. Will read in chunks as I see fit."""
  def __init__(self, saveStream):
    self.civ3 = hexdump(saveStream.read(30),30)
    #print self.civ3

    self.bic = hexdump(saveStream.read(532))
    #print self.bic

    self.bicq = hexdump(saveStream.read(736))
    print self.bicq

    # Plan to read array of arrays until the GAME array is read
    name = ""
    self.bicarray = []
    while name <> "GAME":
      (name, count,) = struct.unpack('4si', saveStream.read(8))
      print name
      print count
      for i in range(count):
        array = []
        (len,) = struct.unpack('i', saveStream.read(4))
        print len
        array.append(hexdump(saveStream.read(len)))
      self.bicarray.append(array)

#    # This isn't always GAME...sometimes it's BLDG
#    #self.game = hexdump(saveStream.read(7593))
#    #print self.game
#
#    # Hmm, this seems to work so far
#    (self.arrayname, self.arraylen,) = struct.unpack('4si', saveStream.read(8))
#    print self.arrayname
#    print self.arraylen
#    self.array = []
#    for i in range(self.arraylen):
#      (len,) = struct.unpack('i', saveStream.read(4))
#      print len
#      self.array.append(hexdump(saveStream.read(len)))
#    print json.dumps(self.array)
#
#    # so let's try it again
#    # Hmm, works on the custom biqs but not the standard ones
#    # Maybe look for GAME and stop
#    (self.array2name, self.array2len,) = struct.unpack('4si', saveStream.read(8))
#    print self.array2name
#    print self.array2len
#    self.array2 = []
#    if self.array2name <> 'GAME':
#      for i in range(self.array2len):
#        (len,) = struct.unpack('i', saveStream.read(4))
#        print len
#        data = hexdump(saveStream.read(len))
#        print data
#        self.array2.append(data)
#    #print json.dumps(self.array2)

    self.whatsnext = hexdump(saveStream.read(40))
    print self.whatsnext


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
    game = Civ3(saveFile, debug)
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

def jsonDefault(o):
    """Trying to make json.dumps() work on all my data"""
    return o.__dict__

def main():
    """If run directly, parse save from argument or stdin and print out hex dumps"""
    saveFile = sys.stdin
    #game = parse_save(saveFile, debug)
    game = newParse(saveFile)
    #print "Output should go here"
    # If this errors out, I probably forgot to delete a buffer from a section
    #print json.dumps(game, default=jsonDefault, indent=4)

if __name__=="__main__":
    main()
