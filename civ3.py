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

# 2014-01-03 I'm in the middle of a major overhaul, and now I'm going to
#  overhaul again. I'm going to try class initiators to read in the
#  file again now that I think I see a more apt class inheritance
#  hierarchy. But I fully realize now there is no "save format" per se.
#  The save is a dump of C++ structures, many of which are descendant
#  classes of a class that writes to disk as a 4-char string followed
#  usually by a data length or an array count. There may or may not be
#  data in the file from other data structures.

import re	# Regular expressions
import struct	# For parsing binary data
import json     # to export JSON
from horspool import horspool    # to seek to first match; from http://inspirated.com/2010/06/19/using-boyer-moore-horspool-algorithm-on-file-streams-in-python
import sys

#class SectionParent:
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
