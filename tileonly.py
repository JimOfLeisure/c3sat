#!/usr/bin/env python

#    Copyright 2013 Jim Nelson
#
#    This file is part of Civ3 Show-And-Tell.
#
#    Civ3 Show-And-Tell is free software: you can redistribute it and/or modify
#    it under the terms of the GNU General Public License as published by
#    the Free Software Foundation, either version 3 of the License, or
#    (at your option) any later version.
#
#    Civ3 Show-And-Tell is distributed in the hope that it will be useful,
#    but WITHOUT ANY WARRANTY; without even the implied warranty of
#    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
#    GNU General Public License for more details.
#
#    You should have received a copy of the GNU General Public License
#    along with Civ3 Show-And-Tell.  If not, see <http://www.gnu.org/licenses/>.

# 2013-04-20 Another attempt at reading a save file; this time I'm going
#   to try extracting just the map data

import struct	# For parsing binary data
import json     # to export JSON for the HTML browser

class GenericSection:
    """Base class for reading SAV sections."""
    def __init__(self, saveStream):
        buffer = saveStream.read(8)
        (self.name, self.length,) = struct.unpack_from('4si', buffer)
        self.offset = saveStream.tell()
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

    def table_out(self):
        """Return a string of a simple text table of visible tiles."""
        table_string = ''
        for y in range(self.height):
            if y % 2 == 1:
                table_string += '  '
            for x in range(self.width / 2):
                if self.tile[x + y * self.width / 2].is_visible:
                    table_string += '#'
                else:
                    table_string += '.'
                table_string += '   '
            table_string += '\n'
        return table_string

    def html_out(self):
        """Return a string of a html table of visible tiles."""
        table_string = '<div class="map">'
        for y in range(self.height):
            table_string += '<div class="maprow">'
            for x in range(self.width / 2):
                i = x + y * self.width /2
                #info = hex(self.tile[i].whatsthis)
                info = str(i)
                if 0 <= i < len(self.tile):
                    if self.tile[i].is_visible_now:
                        table_string += '<div class="tile visible visiblenow">' + info + '</div>'
                    elif self.tile[i].is_visible:
                        table_string += '<div class="tile visible">' + info + '</div>'
                    else:
                        table_string += '<div class="tile fog">' + info + '</div>'
                        pass
                else:
                    table_string += '<div class="tile notile">' + info + '</div>'
            table_string += '</div>'
            table_string += '\n'
        table_string += '</div>'
        return table_string

    def html_fake_iso(self):
        """Return a string of a html table of visible tiles. Filling with blank table entries to position isometric tiles."""
        table_string = '<div class="map">'
        for y in range(self.height):
            table_string += '<div class="maprow">'
            if y % 2 == 1:
                table_string += '<div class="tile notile"></div>'
            for x in range(self.width / 2):
                i = x  + y * self.width /2
                info = hex(self.tile[i].whatsthis)
                #info = str(i)
                if 0 <= i < len(self.tile):
                    if self.tile[i].is_visible_now:
                        table_string += '<div class="tile visible visiblenow">' + info + '</div>'
                    elif self.tile[i].is_visible:
                        table_string += '<div class="tile visible">' + info + '</div>'
                    else:
                        table_string += '<div class="tile fog">' + info + '</div>'
                        pass
                else:
                    table_string += '<div class="tile notile">' + info + '</div>'
                table_string += '<div class="tile notile"></div>'
            table_string += '</div>'
            table_string += '\n'
        table_string += '</div>'
        return table_string

    def isometbrick_out(self):
        """Return a string of a html table of visible tiles. Trying fake-isometric layout with column spanning"""
        table_string = '<table>'
        for y in range(self.height):
            table_string += '<tr>'
            if y % 2 == 1:
                table_string += '<td class="tile notile">.</td>'
            for x in range(self.width / 2):
                i = x  + y * self.width /2
                info = hex(self.tile[i].whatsthis)
                #info = str(i)
                cssclass = 'tile '
                if 0 <= i < len(self.tile):
                    if self.tile[i].is_visible:
                        cssclass += 'visible '
                        if self.tile[i].continent == 6:     # HACK! Hard-coding continent number for ocean on my test save; need to link this to CONT sections in the future
                            cssclass += 'bigblue '
                        table_string += '<td colspan="2" class="' + cssclass + '">' + info + '</td>'
                    else:
                        table_string += '<td colspan="2" class="tile fog">' + info + '</td>'
                #else:
                #    table_string += '<td class="tile notile">' + info + '</td>'
                #table_string += '<td class="tile notile"></td>'
            if y % 2 == 0:
                table_string += '<td class="tile notile">.</td>'
            table_string += '</tr>'
            table_string += '\n'
        table_string += '</table>'
        return table_string

    def json_out(self):
        """Return a string of JSON-formatted tile data."""
        tempmatrix = []
        for y in range(self.height):
            tempmatrix.append([])
            for x in range(self.width / 2):
                tempmatrix[y].append(self.tile_matrix[y][x].info)
                
        #json_string = json.dumps(tempmatrix)
        # Let's try pretty printing and compare uncompressed and compressed sizes
        json_string = json.dumps(tempmatrix, sort_keys=True, indent=4, separators=(',', ': '))
        #json_string = json.JSONEncoder(skipkeys=True).encode(self.tile)
        return json_string

    def map_id(self, i):
        """Return a string to be used as a CSS ID for the tile group. i is the index of self.tile"""
        return 'map' + str(i)

    def svg_out(self):
        """Return a string of svg-coded map"""
        x_axis_wrap = True
        y_axis_wrap = False
        tile_width = 128
        tile_height = 64
        map_width = (self.width * tile_width / 2) + (tile_width / 2)
        map_height = (self.height * tile_height / 2) + (tile_height / 2)
        svg_string = '<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" version="1.1" viewBox="0 0 ' + str(map_width) + ' ' + str(map_height) + '">\n'
        # Well this was dumb. I should just do one big rectangle or change the SVG background
#        if not y_axis_wrap:
#            # Top row border
#            svg_string += '<rect id="topBorder" class="mapEdge" x="0" y="0" width="' + str(map_width) + '" height="' + str(tile_height / 2) + '" />\n'
#            svg_string += '<use xlink:href="#topBorder" transform="translate(0, ' + str(map_height) + ')" />'
#        if not x_axis_wrap:
#            # Left row border
#            svg_string += '<rect id="leftBorder" class="mapEdge" x="0" y="0" width="' + str(tile_width / 2) + '" height="' + str(map_height) + '" />\n'
#            svg_string += '<use xlink:href="#leftBorder" transform="translate(' + str(map_width - (tile_width / 2)) + ', 0)" />'
        svg_string += '<rect class="mapEdge" x="0" y="0" width="' + str(map_width) + '" height="' + str(map_height) + '" />\n'
        for y in range(self.height):
            x_indent = (y % 2) * tile_width / 2
            y_offset = y * tile_height / 2
            svg_string += '<g row="' + str(y) + '" transform="translate(' + str(x_indent) + ', ' + str(y_offset) + ')">\n'
            for x in range(self.width / 2):
                i = x  + y * self.width /2
                info = hex(self.tile[i].whatsthis)
                cssclass = 'tile '
                if 0 <= i < len(self.tile):
                    #svg_string += '  <g transform="translate(' + str((x * tile_width) + x_indent) + ', 0)">\n'
                    svg_string += '  <g id="' + self.map_id(i) + '" transform="translate(' + str(x * tile_width) + ', 0)">\n'
                    # svg_string += '  <g transform="translate(' + str((x * tile_width) + x_indent) + ', ' + str((y % 2) * tile_width / 2) + ')">\n'
                    #svg_string += '    <polygon points="0,20 24,0 48,20 24,40" transform="translate(' + str((x * tile_width) + x_indent) + ', ' + str(y_offset) + ')" '
                    svg_string += '    <polygon points="0,32 64,0 128,32 64,64" '
                    if self.tile[i].is_visible:
                        cssclass += 'visible '
                        if self.tile[i].continent == 6:     # HACK! Hard-coding continent number for ocean on my test save; need to link this to CONT sections in the future
                            cssclass += 'bigblue '
                        svg_string += 'class="' + cssclass + '" />\n'
                        svg_string += '    <text text-anchor="middle" alignment-baseline="central" x="' + str(tile_width / 2) + '" y="' + str(tile_height / 2) + '">' + info + '</text>\n'
                        svg_string += '  </g>\n'
                    else:
                        cssclass = 'tile fog'
                        svg_string += 'class="' + cssclass + '" />\n'
                        ## spoiler info ### svg_string += '    <text >' + info + '</text>\n'
                        svg_string += '  </g>\n'
            # link the first item and place at the end for even rows; link to the last item and place at the first. Will be half-cropped by viewport
            # using math (even lines have 0 remainder, multiplying to cancel out values) instead of if, but it's a little harder to follow
            svg_string += '  <use xlink:href="#' + self.map_id((y * self.width / 2) + (x * (y % 2))) + '" transform="translate(' + str((map_width - tile_width / 2) - (map_width - tile_width / 2) * 2 * (y % 2)) + ', 0)" />\n'
            svg_string += '</g>\n'
        svg_string += '</svg>\n'
        return svg_string

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

def parse_save():
    saveFilePath = "unc-test.sav"
    saveFile = open(saveFilePath, 'rb')
    print 'HACK: Skipping to first TILE in my test SAV.'
    saveFile.seek(0x34a4, 0)
    print 'HACK: Instantiating the class that reads TILEs with width x height hard-coded to my test SAV.'
    game = Tiles(saveFile, 60, 60)
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
    return ''.join(lines)

def main():
    print json
    game = parse_save()
    print 'Printing something(s) from the class to ensure I have what I intended'
    #print game.name, game.length
    #print game.tile.pop()[1].length
    #print game.width, game.height
    #print game.tile[0].Tile128.length
    #print game.tile[1000].Tile128.length
    #print game.tile[0].is_visible_to
    #max = len(game.tile)
    max = 10
    for i in range(max):
        #print i, hex(game.tile[i].Tile36.offset), hex(game.tile[i].whatsthis), hex(game.tile[i].whatsthis2), hex(game.tile[i].whatsthis3), hex(game.tile[i].whatsthis4), hex(game.tile[i].whatsthis5)
        #print i, hex(game.tile[i].Tile36.offset), hex(game.tile[i].whatsthis2)
        #print hex(game.tile[i].whatsthis)
        #print hexdump(game.tile[i].Tile12.buffer)
        print game.tile[i].continent
    #print game.html_out()

if __name__=="__main__":
    main()
