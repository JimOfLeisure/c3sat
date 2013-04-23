#!/usr/bin/env python

# 2013-04-20 Another attempt at reading a save file; this time I'm going
#   to try extracting just the map data

import struct	# For parsing binary data

class GenericSection:
    """Base class for reading SAV sections."""
    def __init__(self, saveStream):
        buffer = saveStream.read(8)
        (self.name, self.length,) = struct.unpack_from('4si', buffer)
        self.buffer = saveStream.read(self.length)

class Tile:
    """Class for each logical tile."""
    def __init__(self, saveStream):

        self.Tile36 = GenericSection(saveStream)
        del self.Tile36.buffer

        self.Tile12 = GenericSection(saveStream)
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
        logical_tiles = width / 2 * height
        while logical_tiles > 0:
            self.tile.append(Tile(saveStream))
            logical_tiles -= 1

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
                if 0 <= i < len(self.tile):
                    if self.tile[i].is_visible:
                        table_string += '<div class="tile visible">' + str(i) + '</div>'
                    else:
                        table_string += '<div class="tile fog">' + str(i) + '</div>'
                        pass
                else:
                    table_string += '<div class="tile notile">' + str(i) + '</div>'
            table_string += '</div>'
            table_string += '\n'
        table_string += '</div>'
        return table_string


def get_int(buffer, offset):
    """Unpack an int from a buffer at the given offest."""
    (the_int,) = struct.unpack('i', buffer[offset:offset+4])
    return the_int

def parse_save():
    saveFilePath = "unc-test.sav"
    saveFile = open(saveFilePath, 'rb')
    print 'HACK: Skipping to first TILE in my test SAV.'
    saveFile.seek(0x34a4, 0)
    print 'HACK: Instantiating the class that reads TILEs with width x height hard-coded to my test SAV.'
    game = Tiles(saveFile, 60, 60)
    return game


def main():
    game = parse_save()
    print 'Printing something(s) from the class to ensure I have what I intended'
    #print game.name, game.length
    #print game.tile.pop()[1].length
    print game.width, game.height
    print game.tile[0].Tile128.length
    print game.tile[1000].Tile128.length
    print game.tile[0].is_visible_to
    i = 0
    max = len(game.tile)
    while i < max:
        print i, game.tile[i].is_visible
        i += 1
    #print game.html_out()

main()
