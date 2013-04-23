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

class Tiles:
    """Class to read all tiles"""
    def __init__(self, saveStream, width, height):
        logical_tiles = width / 2 * height
        self.tile = []
        while logical_tiles > 0:
            Tile36 = GenericSection(saveStream)
            Tile12 = GenericSection(saveStream)
            Tile4 = GenericSection(saveStream)
            Tile128 = GenericSection(saveStream)
            self.tile.append([Tile36, Tile12, Tile4, Tile128])
            logical_tiles -= 1

def parse_save():
    saveFilePath = "unc-test.sav"
    saveFile = open(saveFilePath, 'rb')
    saveFile.seek(0x34a4, 0)
    #game = GenericSection(saveFile)
    game = Tiles(saveFile, 60, 60)
    #print game.name, game.length
    #print game.tile.pop()[1].length
    print game.tile[0][0].length


def main():
    parse_save()

main()
