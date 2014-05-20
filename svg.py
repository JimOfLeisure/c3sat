#!/usr/bin/env python

#    Copyright 2014 Jim Nelson
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

import wrld
#import tileonly
#import datetime
import sys


def main():
    """This module instantiates wrld.parse_save() and writes an svg file for the map"""
    outputsvgpath = 'html/map.svg'
    if len(sys.argv) < 2:
        #print "Usage: svg.py <filename>"
        #sys.exit(-1)
        saveFile = sys.stdin
    else:
        saveFile = open(sys.argv[1], 'rb')

    #game = wrld.parse_save("unc-test.sav")
    #game = wrld.parse_save("unc-lk151-650ad.sav")
    game = wrld.parse_save(saveFile)

    write = open(outputsvgpath, 'w')

    write.write(game.Tiles.svg_out())

main()
