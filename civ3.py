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

# 2015-01-01 Now I want to parse all of any arbitraty game save. Copied from wrld.py. Will probably take spoiler and SVG logic out of this file, too, and put it in the output code

import re	# Regular expressions
import struct	# For parsing binary data
import json     # to export JSON
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
        self.offset = saveStream.tell()
        buffer = saveStream.read(8)
        (self.name, self.length,) = struct.unpack_from('4si', buffer)
        self.buffer = saveStream.read(self.length)
        self.hexdump = hexdump(self.buffer)

class Tile:
    """Class for each logical tile."""
    def __init__(self, saveStream, debug=False):

#        self.info = {}

        self.Tile36 = GenericSection(saveStream)
#        #self.continent = get_short(self.Tile36.buffer, 0x1a)
#        #self.info['continent'] = get_short(self.Tile36.buffer, 0x1a)
#        self.rivers = get_byte(self.Tile36.buffer, 0x00)
#        self.Tile36.values = struct.unpack_from('2h3i2c9h', self.Tile36.buffer)
#        self.continent = self.Tile36.values[11]
#        self.info['continent'] = self.Tile36.values[11]
#        self.top_unit_id = self.Tile36.values[3]
#        self.resource = self.Tile36.values[2]
#        self.barb_info = self.Tile36.values[8]
#        self.city_id = self.Tile36.values[9]
#        self.colony = self.Tile36.values[10]
#        #self.whatsthis = self.Tile4.values[0]
#        self.whatsthis = get_short(self.Tile36.buffer, 0x00) & myglobalbitmask
        if not debug:
#            del self.Tile36.values
            del self.Tile36.buffer

        self.Tile12 = GenericSection(saveStream)
#        self.Tile12.values = struct.unpack_from('iihh', self.Tile12.buffer)
#        self.info['terrain'] = get_byte(self.Tile12.buffer, 0x5)
#        self.improvements = self.Tile12.values[0]
#        self.terrain_features = get_short(self.Tile12.buffer, 0x0a)
#        # Mask 0x0001 is *** Bonus Grassland *** . Interesting, some hills and mountains have it (base tile is grassland)
#        # Mask 0x0002 is the "fat x" around each city whether or not it has the culture to work it
#        # Mask 0x0004 - Can't find any
#        # Mask 0x0008 - Player start location
#        # Mask 0x0010 - Snow-caps for mountains
#        # Mask 0x0020 - Unsure. Only on land, seems clumped. Seems to be all tundra on generated maps and all forest tundra on LK's WM. Have not seen it on jungle tiles.
#        # Other nybble masks 0x00c0 - nothing I can find
#        # Mask 0x1000 - This looks like a likely candidate for "forest already chopped here"
#        # Other nybble masks 0xe000 - nothing I can find ***  (found some of these on LK's WM)
#        # Mask 0x0100 - river N corner? or NW?
#        # Mask 0x0200 - river W corner?
#        # Mask 0x0400 - river E corner? or SE?
#        # Mask 0x0800 - river S corner?
        if not debug:
#            del self.Tile12.values
            del self.Tile12.buffer

        self.Tile4 = GenericSection(saveStream)
#        self.Tile4.values = struct.unpack_from('i', self.Tile4.buffer)
        if not debug:
#            del self.Tile4.values
            del self.Tile4.buffer

        self.Tile128 = GenericSection(saveStream)
#        self.Tile128.values = struct.unpack_from('4i96b', self.Tile128.buffer)
#        #self.is_visible_to = get_int(self.Tile128.buffer, 0)
#        #self.is_visible = self.is_visible_to & 0x02
#        #self.is_visible = self.is_visible_to & 0x10
#        #self.is_visible_now = self.is_visible_now_to & 0x02
#        #self.is_visible_now_to = get_int(self.Tile128.buffer, 4)
#
#        self.is_visible_to_flags = self.Tile128.values[0]
#        self.is_visible_to = []
#        mytemp = self.is_visible_to_flags
#        for civ in range(32):
#            self.is_visible_to.append(mytemp & 0x01 == 1)
#            mytemp = mytemp >> 1
#
#        self.is_visible_now_to_flags = self.Tile128.values[1]
#        self.is_visible_now_to = []
#        mytemp = self.is_visible_now_to_flags
#        for civ in range(32):
#            self.is_visible_now_to.append(mytemp & 0x01 == 1)
#            mytemp = mytemp >> 1
#
#        self.worked_by_city_id = get_short(self.Tile128.buffer, 0x14)
#
#        mytemp = 0x16
#        self.land_trade_network_id = []
#        for civ in range(32):
#            self.land_trade_network_id.append(get_short(self.Tile128.buffer, mytemp))
#            mytemp += 2
#
#        mytemp = 0x56
#        self.improvements_known_to_civ = []
#        for civ in range(32):
#            self.improvements_known_to_civ.append(get_byte(self.Tile128.buffer, mytemp))
#            mytemp += 1
#
        if not debug:
#            del self.Tile128.values
            del self.Tile128.buffer


class Tiles:
    """Class to read all tiles"""
    def __init__(self, saveStream, width, height, debug=False):
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
                this_tile = Tile(saveStream, debug)

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
          myriver = self.tile[i].rivers
          if myriver <> 0:
		      if myriver & 0x80 <> 0: mystring += '<use xlink:href="#floodplain-nw" ' + xypos +' />\n'
		      if myriver & 0x02 <> 0: mystring += '<use xlink:href="#floodplain-ne" ' + xypos +' />\n'
		      if myriver & 0x08 <> 0: mystring += '<use xlink:href="#floodplain-se" ' + xypos +' />\n'
		      if myriver & 0x20 <> 0: mystring += '<use xlink:href="#floodplain-sw" ' + xypos +' />\n'
		      if myriver == 0x01: mystring += '<use xlink:href="#floodplain-n-corner" ' + xypos +' />\n'
		      if myriver == 0x10: mystring += '<use xlink:href="#floodplain-s-corner" ' + xypos +' />\n'
		      if myriver == 0x04: mystring += '<use xlink:href="#floodplain-e-corner" ' + xypos +' />\n'
		      if myriver == 0x40: mystring += '<use xlink:href="#floodplain-w-corner" ' + xypos +' />\n'
        elif base_terrain == 1: mystring = '<use xlink:href="#plains" ' + xypos +' />\n'
        elif base_terrain == 2: mystring = '<use xlink:href="#grassland" ' + xypos +' />\n'
        elif base_terrain == 3: mystring = '<use xlink:href="#tundra" ' + xypos +' />\n'
        elif base_terrain == 11: mystring = '<use xlink:href="#coast" ' + xypos +' />\n'
        elif base_terrain == 12: mystring = '<use xlink:href="#sea" ' + xypos +' />\n'
        elif base_terrain == 13: mystring = '<use xlink:href="#ocean" ' + xypos +' />\n'
        else: mystring = '<use xlink:href="#unknown" ' + xypos +' />\n'
        return mystring

    def andeq(self,a,b):
        return a & b == b

    def rivers(self, i, x, y):
      xypos = self.svg_attr_xy((x,y))
      myriver = self.tile[i].rivers
      mystring = ""
      if myriver & 0x80 <> 0: mystring += '<use xlink:href="#river-nw" ' + xypos +' />\n'
      if myriver & 0x02 <> 0: mystring += '<use xlink:href="#river-ne" ' + xypos +' />\n'
      if myriver & 0x08 <> 0: mystring += '<use xlink:href="#river-se" ' + xypos +' />\n'
      if myriver & 0x20 <> 0: mystring += '<use xlink:href="#river-sw" ' + xypos +' />\n'
      if myriver == 0x01: mystring += '<use xlink:href="#river-n-corner" ' + xypos +' />\n'
#      if self.andeq(myriver,0xc0): mystring += '<use xlink:href="#river-nw" ' + xypos +' />\n'
#      if self.andeq(myriver,0x06): mystring += '<use xlink:href="#river-ne" ' + xypos +' />\n'
#      if self.andeq(myriver,0x1c): mystring += '<use xlink:href="#river-se" ' + xypos +' />\n'
#      if self.andeq(myriver,0x20): mystring += '<use xlink:href="#river-sw" ' + xypos +' />\n'
      return mystring

    def overlay_terrain(self, i, x, y):
      xypos = self.svg_attr_xy((x,y))
      textxypos = self.svg_attr_xy((x + self.tile_width /2,y + self.tile_height /2))
      # Get left-nibble of terrain byte: bit-rotate right 4, then mask to be sure it wasn't more than a byte
      overlay_terrain = (self.tile[i].info['terrain'] >> 4) & 0x0F
      if overlay_terrain == 0x04:
          # Flood plain
          mystring = self.svg_text("FP",textxypos)
          mystring = ""

      elif overlay_terrain == 0x05: mystring = '<use ' + xypos + ' xlink:href = "#myHill" />\n'
      elif overlay_terrain == 0x06: mystring = '<use ' + xypos + ' xlink:href = "#myMountain" />\n'
      elif overlay_terrain == 0x07: mystring = '<use ' + xypos + ' xlink:href = "#myForest" />\n'
      elif overlay_terrain == 0x08: mystring = '<use ' + xypos + ' xlink:href = "#myJungle" />\n'
      elif overlay_terrain == 0x09: mystring = '<use ' + xypos + ' xlink:href = "#myMarsh" />\n'
      elif overlay_terrain == 0x0a: mystring = '<use ' + xypos + ' xlink:href = "#myVolcano" />\n'
      elif overlay_terrain == 0x02:
          if self.tile[i].terrain_features & 0x0001 == 1: mystring = '<use xlink:href="#bonusgrassland" ' + xypos +' />\n'
          else: mystring = ""
      elif overlay_terrain in {0,1,2,3,0xb,0xc,0xd}:
          # It appears if there is no overlay, the nybble matches the base tile nybble. Return nothing for known base tile values
          mystring = ""
      else:
          # If something unexpected, put the nybble value in text on the map
          mystring = self.svg_text("0x%01x" % overlay_terrain,textxypos)
      return mystring

    def resource(self, i, x, y):
      xypos = self.svg_attr_xy((x,y))
      textxypos = self.svg_attr_xy((x + self.tile_width /2,y + self.tile_height /2))
      resource = self.tile[i].resource
      if resource == 0x08: mystring = '<use ' + xypos + ' xlink:href = "#myWines" />\n'
      elif resource == 0x09: mystring = self.svg_text("Furs",textxypos)
      elif resource == 0x0c: mystring = self.svg_text("Spices",textxypos)
      elif resource == 0x0d: mystring = self.svg_text("Ivory",textxypos)
      elif resource == 0x0a: mystring = self.svg_text("Dyes",textxypos)
      elif resource == 0x0b: mystring = self.svg_text("Incense",textxypos)
      elif resource == 0x0e: mystring = self.svg_text("Silks",textxypos)
      elif resource == 0x0f: mystring = self.svg_text("Gems",textxypos)
      elif resource == 0x10: mystring = self.svg_text("Whales",textxypos)
      elif resource == 0x11: mystring = self.svg_text("Deer",textxypos)
      elif resource == 0x12: mystring = '<use ' + xypos + ' xlink:href = "#myFish" />\n'
      #elif resource == 0x13: mystring = self.svg_text("Cow",textxypos)
      elif resource == 0x13: mystring = '<use ' + xypos + ' xlink:href = "#myCow" />\n'
      elif resource == 0x14: mystring = self.svg_text("Wheat",textxypos)
      elif resource == 0x15: mystring = self.svg_text("Gold",textxypos)
      elif resource == 0x16: mystring = self.svg_text("Sugar",textxypos)
      elif resource == 0x17: mystring = self.svg_text("Trop.Fruit",textxypos)
      elif resource == 0x18: mystring = self.svg_text("Oasis",textxypos)
      elif resource == 0x19: mystring = self.svg_text("Tobacco",textxypos)
      elif resource in {-1,0,1,2,3,4,5,6,7}:
          # Treating all strategic resources as spoiler info
          # Also no string for -1, no resource
          mystring = ""
      else:
          # If something unexpected, put the value in text on the map
          mystring = self.svg_text("0x%01x" % resource,textxypos)
      return mystring

    def city(self, i, x, y):
      xypos = self.svg_attr_xy((x,y))
      textxypos = self.svg_attr_xy((x + self.tile_width /2,y + self.tile_height /2))
      city = self.tile[i].city_id
      if city <> -1:
        mystring = '<use ' + xypos + ' xlink:href = "#myCity" />\n'
      else:
        mystring = ""
      return mystring

    def debug_text(self,i,x,y):
        xypos = self.svg_attr_xy((x + self.tile_width /2,y + self.tile_height /2))
        (mapx,mapy) = self.map_xy(i)

        # x,y in hex:
        #mystring = self.svg_text("%02x" % mapx + ',' + "%02x" % mapy, xypos)

        # "(x,y)" in decimal:
        #mystring = self.svg_text(str(self.map_xy(i)), xypos)

        # i in decimal:
        #mystring = self.svg_text(str(i), xypos)

        # continent ID
        #mystring = self.svg_text(str(self.tile[i].continent), xypos)

        #
        #subject = self.tile[i].Tile36.values[6]
        subject = self.tile[i].whatsthis
        #if subject <> -1:
        #if subject <> 0:
        if subject <> myglobalcompare:
        #if subject <> self.tile[i].improvements:
        #if subject & 0x20 <> 0:
            mystring = self.svg_text("%04x" % subject, xypos)
        else:
            mystring = ""


        return mystring

    def svg_out(self, spoiler=False, debug=False):
        """Return a string of svg-coded map"""
        x_axis_wrap = True
        y_axis_wrap = False
        self.tile_width = 128
        self.tile_height = 64
        map_width = (self.width * self.tile_width / 2) + (self.tile_width / 2)
        map_height = (self.height * self.tile_height / 2) + (self.tile_height / 2)
        svg_string = ""
        svg_string += '<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" version="1.1" x="0" y="0" viewBox="0 0 ' + str(map_width) + ' ' + str(map_height) + '">\n'
        svg_string += "<defs>\n"
        mapDefsFile = open("mapdefs.svg","r")
        svg_string += mapDefsFile.read()
        mapDefsFile.close()
        #svg_string += '<g id="myCow" transform="scale(0.12) translate(320,117)">\n'
        #mapDefsFile = open("svg/Cow_cartoon_04.svg","r")
        #svg_string += mapDefsFile.read()
        #mapDefsFile.close()
        #svg_string += '</g>\n'
        svg_string += "</defs>\n"
        svg_string += '<use xlink:href="#mybackgroundrectangle" x="0" y="0" transform="scale(' + str(map_width) + ',' + str(map_height) + ')" />\n'
        for i in range(len(self.tile)):
          #if self.tile[i].is_visible_to[myglobalciv] or spoiler:
          if self.tile[i].is_visible_to[myglobalciv] or spoiler:
            # May have more than one to paint wrap-around tiles
            for (x,y) in self.svg_xy(i):
              svg_string += self.base_terrain(i,x,y)
              svg_string += self.rivers(i,x,y)
              svg_string += self.overlay_terrain(i,x,y)
              svg_string += self.resource(i,x,y)
              svg_string += self.city(i,x,y)
              #if debug: svg_string += self.svg_text(str(self.map_xy(i)),self.svg_attr_xy((x + self.tile_width /2,y + self.tile_height /2)))
              if debug: svg_string += self.debug_text(i,x,y)
          #else: # I used to place a fog tile, but why not just let the background rectangle be the fog?
          #  for (x,y) in self.svg_xy(i):
          #    svg_string += '<use xlink:href="#fog" ' + self.svg_attr_xy((x,y)) +' />\n'
        if spoiler:
            for civ in range(len(self.start_loc)):
                i = self.start_loc[civ]
                if i <> -1:
                    for (x,y) in self.svg_xy(i):
                        xypos = self.svg_attr_xy((x + self.tile_width /2,y + self.tile_height /2))
                        svg_string += self.svg_text('Player ' + str(civ) + ' Start', xypos)
        svg_string += '</svg>\n'
        return svg_string

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
        self.name = saveStream.read(4)
        # Reading an arbitray # of bytes until I figure out the format
        self.buffer = saveStream.read(400)
        self.hexdump = hexdump(self.buffer)
        if not debug:
          del self.buffer

class Civ3:
    """Class for entire save game"""
    def __init__(self, saveStream, debug=False):
        self.name = saveStream.read(4)
        if self.name <> 'CIV3':
            print "wah wah wah wahhhhhhhh."
            print "Stub. Provided stream not decompressed C3C save"
            return -1
        self.buffer = saveStream.read(26)
        self.hexdump = hexdump(self.buffer, 26)
        if not debug:
          del self.buffer
        self.Bic = Bic(saveStream, debug)
        #print 'Using Horspool search to go to first WRLD section'
        wrldOffset = horspool.boyermoore_horspool(saveStream, "WRLD")
        #print wrldOffset
        self.Wrld = Wrld(saveStream, debug)
        saveStream.close()


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
    spoiler = True
    spoiler = False
    debug = True
    debug = False
    if len(sys.argv) < 2:
        #print "Usage: svg.py <filename>"
        #sys.exit(-1)
        saveFile = sys.stdin
    else:
        saveFile = open(sys.argv[1], 'rb')
    game = parse_save(saveFile, debug)
    print "Output should go here"
    print json.dumps(game, default=jsonDefault, indent=4)
    # Maybe I should just dump new things since I have WRLD and TILE more or less figured out

if __name__=="__main__":
    main()