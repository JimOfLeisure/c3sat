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
#import json     # to export JSON for the HTML browser
import horspool    # to seek to first match; from http://inspirated.com/2010/06/19/using-boyer-moore-horspool-algorithm-on-file-streams-in-python
import sys

class GenericSection:
    """Base class for reading SAV sections."""
    def __init__(self, saveStream):
        buffer = saveStream.read(8)
        (self.name, self.length,) = struct.unpack_from('4si', buffer)
        #self.offset = saveStream.tell()
        self.buffer = saveStream.read(self.length)

class Tile:
    """Class for each logical tile."""
    def __init__(self, saveStream):

        self.info = {}

        self.Tile36 = GenericSection(saveStream)
        self.continent = get_short(self.Tile36.buffer, 0x1a)
        self.info['continent'] = get_short(self.Tile36.buffer, 0x1a)
        del self.Tile36.buffer

        self.Tile12 = GenericSection(saveStream)
        self.info['terrain'] = get_byte(self.Tile12.buffer, 0x5)
        #self.whatsthis = get_byte(self.Tile12.buffer, 0xa)
        self.whatsthis = get_byte(self.Tile12.buffer, 0x5)
        self.whatsthis2 = get_byte(self.Tile12.buffer, 0x5)
        self.whatsthis3 = get_byte(self.Tile12.buffer, 0xa)
        self.whatsthis4 = get_byte(self.Tile12.buffer, 0xb)
        self.whatsthis5 = get_int(self.Tile12.buffer, 0x0)
        del self.Tile12.buffer

        self.Tile4 = GenericSection(saveStream)
        del self.Tile4.buffer

        self.Tile128 = GenericSection(saveStream)
        self.is_visible_to = get_int(self.Tile128.buffer, 0)
        self.is_visible_now_to = get_int(self.Tile128.buffer, 4)
        self.is_visible = self.is_visible_to & 0x02
        #self.is_visible = self.is_visible_to & 0x10
        self.is_visible_now = self.is_visible_now_to & 0x02
        del self.Tile128.buffer

class Tiles:
    """Class to read all tiles"""
    def __init__(self, saveStream, width, height):
        self.width = width      # These may eventually be redundant to a parent class
        self.height = height
        self.tile = []          # List of individual tiles
        self.tile_matrix = []       # x,y matrix of individual tiles
        self.tile_iso_matrix = []   # faux isometric padded matrix  of individual tiles
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
                this_tile = Tile(saveStream)

                self.tile.append(this_tile)

                self.tile_matrix[y].append(this_tile)

                self.tile_iso_matrix[y].append(this_tile)
                self.tile_iso_matrix[y].append(None)

            if y % 2 == 0:
                self.tile_iso_matrix[y].append(None)

#    def map_id(self, i):
#        """Return a string to be used as a CSS ID for the tile group. i is the index of self.tile"""
#        return 'map' + str(i)

    def map_xy(self, i):
        """Return logical fake-isometric x,y coordinates for tile[i]"""
        # Due to fake isometric layout, x+2 is one tile East, and there are mapwidth / 2 tiles per row
        y = i // (self.width // 2)
        # Odd rows are offset
        x = 2 * (i % (self.width // 2)) + (y % 2)
        return (x,y)

    def svg_xy(self, i):
        """Return x,y position on SVG canvas for tile[i]. May be more than one result to allow for edge tile wrapping"""
        (x,y) = self.map_xy(i)
        svg_x = x * self.tile_width / 2
        svg_y = y * self.tile_height / 2
        result = ((svg_x,svg_y),)
        # Add SVG coords for wraparound tiles where appropriate
        if y % 2 == 0:
            if x == 0:
                result += (((self.width / 2) * self.tile_width, svg_y),)
        else:
            if x + 1 == self.width:
                result += ((0 - self.tile_width / 2,svg_y),)

        return result

    def svg_attr_xy(self, (x, y)):
        """Return string for x and y attributes for SVG shapes for a given tile"""
        return 'x="' + str(x) + '" y="' + str(y) + '"'

    def test_things(self):
        """Test function to develop x/y functions"""
        self.tile_width = 128
        self.tile_height = 64
        for i in range(180):
            print self.map_xy(i)
            print "     ",  self.svg_xy(i)
            for (x,y) in self.svg_xy(i):
                print 'x="' + str(x) + '" y="' + str(y) + '"'
            #print self.svg_attr_xy(self.svg_xy(i))
        return "Yay"

    def svg_text(self, text, xypos):
        return '<text ' + xypos + ' text-anchor="middle" alignment-baseline="central" style="font-size:32px">' + str(text) + '</text>\n'
        #return '<text ' + xypos + '>' + str(text) + '</text>\n'

    def base_terrain(self, i, x, y):
        """Return a string for base terrain"""
        xypos = self.svg_attr_xy((x,y))
        # Get right-nibble of terrain byte
        base_terrain = self.tile[i].info['terrain'] & 0x0F
        if base_terrain == 0:
            mystring = '<use xlink:href="#desert" ' + xypos +' />\n'
        elif base_terrain == 1:
            mystring = '<use xlink:href="#plains" ' + xypos +' />\n'
        elif base_terrain == 2:
            mystring = '<use xlink:href="#grassland" ' + xypos +' />\n'
        elif base_terrain == 3:
            mystring = '<use xlink:href="#tundra" ' + xypos +' />\n'
        elif base_terrain == 11:
            mystring = '<use xlink:href="#coast" ' + xypos +' />\n'
        elif base_terrain == 12:
            mystring = '<use xlink:href="#sea" ' + xypos +' />\n'
        elif base_terrain == 13:
            mystring = '<use xlink:href="#ocean" ' + xypos +' />\n'
        else:
            mystring = '<use xlink:href="#unknown" ' + xypos +' />\n'
        return mystring

    def overlay_terrain(self, i, x, y):
      xypos = self.svg_attr_xy((x,y))
      textxypos = self.svg_attr_xy((x + self.tile_width /2,y + self.tile_height /2))
      # Get left-nibble of terrain byte: bit-rotate right 4, then mask to be sure it wasn't more than a byte
      overlay_terrain = (self.tile[i].info['terrain'] >> 4) & 0x0F
      if overlay_terrain == 0x04:
          # Flood plain
          #mystring = '    <text class="' + cssclass + '" text-anchor="middle" alignment-baseline="central" x="' + xypos + '">Flood Plain</text>\n'
          mystring = self.svg_text("Flood Plain",textxypos)
      elif overlay_terrain == 0x05:
          # Hill
          mystring = '    <use ' + xypos + ' xlink:href = "#myHill" />\n'
      elif overlay_terrain == 0x06:
          # Mountain
          mystring = '    <use ' + xypos + ' xlink:href = "#myMountain" />\n'
      elif overlay_terrain == 0x07:
          # Forest
          mystring = '    <use ' + xypos + ' xlink:href = "#myForest" />\n'
      elif overlay_terrain == 0x08:
          # Jungle
          #mystring = '    <text class="' + cssclass + '" text-anchor="middle" alignment-baseline="central" x="' + xypos + '">Jungle</text>\n'
          mystring = self.svg_text("Jungle",textxypos)
      elif overlay_terrain == 0x09:
          # Marsh
          #mystring = '    <text class="' + cssclass + '" text-anchor="middle" alignment-baseline="central" x="' + xypos + '">Marsh</text>\n'
          mystring = self.svg_text("Marsh",textxypos)
      elif overlay_terrain == 0x0a:
          # Volcano
          mystring = '    <use ' + xypos + ' xlink:href = "#myVolcano" />\n'
      elif overlay_terrain in {0,1,2,3,11,12,13}:
          mystring = ""
      else:
          #mystring = '    <text class="' + cssclass + '" text-anchor="middle" alignment-baseline="central" x="' + xypos + '">' + info + '</text>\n'
          #mystring = self.svg_text(overlay_terrain,xypos)
          mystring = self.svg_text("0x%01x" % overlay_terrain,textxypos)
      return mystring

    def svg_out(self, spoiler=False):
        """Return a string of svg-coded map"""
        x_axis_wrap = True
        y_axis_wrap = False
        self.tile_width = 128
        self.tile_height = 64
        map_width = (self.width * self.tile_width / 2) + (self.tile_width / 2)
        map_height = (self.height * self.tile_height / 2) + (self.tile_height / 2)
        svg_string = ""
        svg_string += '<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" version="1.1" x="0" y="0" viewBox="0 0 ' + str(map_width) + ' ' + str(map_height) + '">\n'
        #svg_string += '<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" version="1.1" width="' + str(map_width) + '" height="' + str(map_height) + '">\n'
        svg_string += "<defs>\n"
        mapDefsFile = open("mapdefs.svg","r")
        svg_string += mapDefsFile.read()
        mapDefsFile.close()
        svg_string += "</defs>\n"
        svg_string += '<use xlink:href="#mybackgroundrectangle" x="0" y="0" transform="scale(' + str(map_width) + ',' + str(map_height) + ')" />\n'
        for i in range(len(self.tile)):
          if self.tile[i].is_visible or spoiler:
            for (x,y) in self.svg_xy(i):
              #xypos = self.svg_attr_xy((x,y))
              svg_string += self.base_terrain(i,x,y)
              svg_string += self.overlay_terrain(i,x,y)
#        for y in range(self.height):
#            x_indent = (y % 2) * tile_width / 2
#            y_offset = y * tile_height / 2
#            #svg_string += '<g row="' + str(y) + '" transform="translate(' + str(x_indent) + ', ' + str(y_offset) + ')">\n'
#            for x in range(self.width / 2):
#                i = x  + y * self.width /2
#                info = hex(self.tile[i].whatsthis)
#                cssclass = 'tile'
#                if 0 <= i < len(self.tile):
#                    #svg_string += '  <g id="' + self.map_id(i) + '" transform="translate(' + str(x * tile_width) + ', 0)">\n'

#                        # Not sure I need this comparison; may just be able to key off the nibble value
#                        if overlay_terrain <> base_terrain:
#                            cssclass = 'overlayterrain terroverlay' + str(overlay_terrain)
#                        svg_string += '    <text class="whatsthis" style="display:none" text-anchor="middle" alignment-baseline="central" x="' + str(tile_width / 2) + '" y="' + str(tile_height / 2) + '">' + info + '</text>\n'
#                        #svg_string += '  </g>\n'
#                    else:
#                        cssclass = 'fog'
#                        svg_string += '<use xlink:href="#fog" class="' + cssclass +'" />\n'
#                        #svg_string += '  </g>\n'
#                    svg_string += '  </g>\n'
#            # link the first item and place at the end for even rows; link to the last item and place at the first. Will be half-cropped by viewport
#            # using math (even lines have 0 remainder, multiplying to cancel out values) instead of if, but it's a little harder to follow
#            svg_string += '  <use xlink:href="#' + self.map_id((y * self.width / 2) + (x * (y % 2))) + '" transform="translate(' + str((map_width - tile_width / 2) - (map_width - tile_width / 2) * 2 * (y % 2)) + ', 0)" />\n'
#            svg_string += '</g>\n'
        svg_string += '</svg>\n'
        return svg_string

class Wrld:
    """Class for 3 WRLD sections"""
    def __init__(self, saveStream):
        """Currently calling this from the horspool seek, so WRLD is already consumed from the stream. Read the length first."""
        self.name = "WRLD"
        buffer = saveStream.read(4)
        (self.length,) = struct.unpack_from('i', buffer)
        #print self.length
        self.buffer = saveStream.read(self.length)
        # Extract any data here, but I think it's only 2 bytes
        #print self.name
        #print hexdump(self.buffer)
        del self.buffer

        self.Wrld2 = GenericSection(saveStream)
        #print self.Wrld2.name
        self.values = struct.unpack_from('41i', self.Wrld2.buffer)
        self.height = self.values[1]
        self.width = self.values[6]
        #print self.height
        #print self.width
        #print self.values
        #print hexdump(self.Wrld2.buffer)
        del self.Wrld2.buffer

        self.Wrld3 = GenericSection(saveStream)
        #print self.Wrld3.name
        #print hexdump(self.Wrld3.buffer)
        del self.Wrld3.buffer

        self.Tiles = Tiles(saveStream, self.width, self.height)


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

#FAIL
#def decompress(firstbytes, saveStream):
#    """Decompress the presumed save file stream. First 4 bytes are already consumed, so we take those in as firstbytes parameter"""
#    #inputStream = StringIO.StringIO()
#    #outputStream = StringIO.StringIO()
#    #inputStream, outputStream = os.pipe()
#    #process = subprocess.Popen(['./blast'], stdin=inputStream,stdout=outputStream, shell=True)
#    process = subprocess.Popen(['./blast'], stdin=subprocess.PIPE,stdout=saveStream, shell=True)
#    process.communicate(firstbytes)
#    process.communicate(saveStream)
#    #response = process.communicate(firstbytes)
#    #response += process.communicate(saveStream)
#    #return outputStream

def parse_save(saveFile):
    buffer = saveFile.read(4)
    if buffer <> 'CIV3':
        print "wah wah wah wahhhhhhhh."
        print "Stub. Provided stream not decompressed C3C save"
        return -1
    #print 'Using Horspool search to go to first WRLD section'
    wrldOffset = horspool.boyermoore_horspool(saveFile, "WRLD")
    #print wrldOffset
    game = Wrld(saveFile)
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
    game.Tiles.test_things()

if __name__=="__main__":
    main()
