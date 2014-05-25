#!/usr/bin/env python

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

import wrld
#import tileonly
#import datetime
import sys


def main():
    """This module instantiates wrld.parse_save() and writes an svg file for the map"""
    outputsvgpath = 'html/civmap.svg'
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

    #write.write(game.Tiles.svg_out(True))
    write.write(game.Tiles.svg_out())

main()
