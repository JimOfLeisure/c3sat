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

# Found useful info at http://stackoverflow.com/questions/192109/is-there-a-function-in-python-to-print-all-the-current-properties-and-values-of
#    and http://stackoverflow.com/questions/109087/python-get-instance-variables

import wrld
from pprint import pprint
import struct


def main():
    """This module instantiates wrld.parse_save() and performs various hex dumps to help me figure out the save format"""
    game = ()
    spoiler = False
    debug = True
    for name in ["test","test1020ad","lk151-650ad","got-map"]:
#    for name in ["test","test1020ad","iso4k","iso170"]:
        path_name = "gamesaves/unc-" + name + ".sav"
        save_file = open(path_name, 'rb')
        print path_name
        game = wrld.parse_save(save_file,debug)
        save_file.close()
#        out_name = "html/debug/idec-" + name + ".html"
#        svg_out = open(out_name, 'w')
#        infile = open("html/debug/head.html","r")
#        svg_out.write(infile.read())
#        infile.close()
#        svg_out.write(game.Tiles.svg_out(spoiler,debug))
#        infile = open("html/debug/tail.html","r")
#        svg_out.write(infile.read())
#        svg_out.close()
        subject = game.Tiles.tile[685].Tile4
        values = struct.unpack_from("i", subject.buffer)
        print wrld.hexdump(subject.buffer);
        print values
        #print game.Tiles.tile[621].improvements_known_to_civ
        #pprint (vars(subject.is_visible_to))
        #pprint (vars(subject))

#    for n in range(len(game)):
        #print game[n].__dict__
        #print game[n].__dict__.keys()
        #pprint (vars(game[n]))
#        print n


main()
