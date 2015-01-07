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

class Section:
    """Top-level class for reading sections / serialized C++ class dumps from safe file"""
    def __init__(self, saveStream, expectedName = None, length = None):
        # Apparently tell() doesn't work on stdin / pipes
        #self.offset = saveStream.tell()
        self.expectedName = expectedName
        self.expectedLength = length
        self.noexist = None
        self.readHeader(saveStream)
        self.readData(saveStream)

    def readHeader(self, saveStream):
        (self.name,) = struct.unpack('4s', saveStream.read(4))
        self.length = self.expectedLength

    def readData(self, saveStream):
        if self.length:
            self.data = saveStream.read(self.length)

    def dumpHeader(self):
        return 'Expected: Name: {0}  Length: {1}\nActual:   Name: {2}  Length: {3}\n'.format(self.expectedName, self.expectedLength, self.name, self.length)

    def dumpData(self):
        return '{0}\n'.format(hexdump(self.data))

    def dumpSelf(self):
        return self.dumpHeader() + self.dumpData()

    def __str__(self):
        return self.dumpSelf()

class NameLength(Section):
    """A typical record that starts with a 4-char sequence followed by length of data in bytes"""
    def readHeader(self, saveStream):
        (self.name, self.length) = struct.unpack('4si', saveStream.read(8))

class HorspoolNameLength(NameLength):
    """A NameLength record where the name has been consumed by a horspool seek"""
    def readHeader(self, saveStream):
        (self.length,) = struct.unpack('i', saveStream.read(4))
        self.name = self.expectedName

class ObjectArray(NameLength):
    """In this type of record the 4-char sequence is followed by a length in records, each of which starts with an integer length in bytes"""
    def readData(self, saveStream):
        self.data = []
        for i in range(self.length):
            (length,) =  struct.unpack('i', saveStream.read(4))
            data = saveStream.read(length)
            self.data.append(data)
    def dumpData(self):
        outtext = []
        for i in range(len(self.data)):
            outtext.append(hexdump(self.data[i]))
        separator = "\n" + "-" * 53 + "\n"
        return separator.join(outtext)

class Flavor():
    """The data unit of FLAV"""
    def __init__(self, saveStream):
        (self.number, self.name, self.numRecords) =   struct.unpack('i256si', saveStream.read(264))
        self.records = []
        for i in range(self.numRecords):
            (someinteger,) = struct.unpack('i', saveStream.read(4))
            self.records.append(someinteger)

class FlavSection(NameLength):
    """The FLAV section of the BIQ is different"""
    def readData(self, saveStream):
        self.flavorGroups = []
        for i in range(self.length):
            (flavors,) =  struct.unpack('i', saveStream.read(4))
            flavor = []
            for j in range(flavors):
                flavor.append(Flavor(saveStream))
        self.flavorGroups.append(flavor)
        self.data = "dummy data for the FLAV class. TODO: FIXME"
#    def dumpData(self):
#        outtext = []
#        for i in range(len(self.flavorGroups)):
#            for j in range(len(self.flavorGroups[i])):
#                outtext.append(Flavor.__str__(self.flavorGroups[i][j]))
#        separator = "\n" + "-" * 53 + "\n"
#        return separator.join(outtext)

class Bic(ObjectArray):
    """The embedded BIC is a little more complex and different, but mostly object arrays"""
    def __init__(self, saveStream, expectedName = None, length = None):
        ObjectArray.__init__(self, saveStream, expectedName, length)
        self.defs = []
        self.defs.append(ObjectArray(saveStream))
        # If GAME is not the first section, then the other rules appear (? Surely there is a flag or other indictor for this content somewhere)
        if self.defs[0].name == "BLDG":
            # The order is assumed to always be the same, but I'm not sure of this yet
            self.defs.append(ObjectArray(saveStream, "CTZN"))
            self.defs.append(ObjectArray(saveStream, "CULT"))
            self.defs.append(ObjectArray(saveStream, "DIFF"))
            self.defs.append(ObjectArray(saveStream, "ERAS"))
            self.defs.append(ObjectArray(saveStream, "ESPN"))
            self.defs.append(ObjectArray(saveStream, "EXPR"))
            # NOTE: In the acutal BIQ file, FLAV seems to come after WSIZ
            self.defs.append(FlavSection(saveStream))
            self.defs.append(ObjectArray(saveStream, "GOOD"))
            self.defs.append(ObjectArray(saveStream, "GOVT"))
            self.defs.append(ObjectArray(saveStream, "RULE"))
            self.defs.append(ObjectArray(saveStream, "PRTO"))
            self.defs.append(ObjectArray(saveStream, "RACE"))
            self.defs.append(ObjectArray(saveStream, "TECH"))
            self.defs.append(ObjectArray(saveStream, "TRFM"))
            self.defs.append(ObjectArray(saveStream, "TERR"))
            self.defs.append(ObjectArray(saveStream, "WSIZ"))
            self.defs.append(ObjectArray(saveStream, "GAME"))
            self.defs.append(ObjectArray(saveStream, "LEAD"))
    def readHeader(self, saveStream):
        (self.name, self.verNum, self.length) = struct.unpack('4s4si', saveStream.read(12))
    def dumpData(self):
        outtext = []
        outtext.append(ObjectArray.dumpData(self))
        for i in range(len(self.defs)):
            outtext.append(self.defs[i].dumpSelf())
        separator = "\n" + "-" * 53 + "\n"
        return separator.join(outtext)

class Idls():
    """The IDLS section seems unique if somewhat like FLAV"""
    def __init__(self, saveStream):
        (self.name, self.someNumber, self.numRecords) =   struct.unpack('4sii', saveStream.read(12))
        self.records = []
        for i in range(self.numRecords):
            (numIntegers,) = struct.unpack('i', saveStream.read(4))
            integers = []
            for j in range(numIntegers):
                (integer,) = struct.unpack('i', saveStream.read(4))
                integers.append(integer)
            self.records.append(integers)

class Popd(ObjectArray):
    """Popd appears to be an object array with an extra number in front of the length"""
    def readHeader(self, saveStream):
        (self.name, self.integer, self.integer2, self.length,) = struct.unpack('4siii', saveStream.read(16))
        print self.name, self.integer, self.integer2, self.length
    def readData(self, saveStream):
        self.data = []
        for i in range(self.length):
            self.data.append(NameLength(saveStream, 'CTZN'))
#            print self.data[-1].name
#            print self.data[-1].length
#            print hexdump(self.data[-1].data)

class newParse:
    """Starting over with parsing strategy. Will read in chunks as I see fit."""
    def __init__(self, saveStream):
        self.civ3 = Section(saveStream, 'CIV3', 26)
        self.bic = NameLength(saveStream, 'BIC ', 524)
        self.embeddedBic = Bic(saveStream, 'BICQ', 1)

        # GAME
        # Trying to  figure it out based on https://github.com/Antal1987/C3CPatchFramework/blob/master/Civ3/Game.h
        # Looks like a bunch of integers
        # SKIPPING over GAME section since I haven't figured it out yet
        #self.gameLength = horspool.boyermoore_horspool(saveStream, "DATE")
        #print 'GAME section length: {0}'.format(self.gameLength)

        #self.game = struct.unpack('46i', saveStream.read(46*4))
        #self.game = struct.unpack('74i', saveStream.read(74*4))
        #self.game = struct.unpack('96i', saveStream.read(96*4))
        #self.game = struct.unpack('263i', saveStream.read(263*4))
        #print len(self.game)
        #print self.game

        print "\nWhat's Next:"
        self.whatsnext = hexdump(saveStream.read(263*4))
        print self.whatsnext

        ### EXITING EARLY WHILE DECODING GAME
        return

        self.date1 = HorspoolNameLength(saveStream, 'DATE', 84)
        self.plgi1 = NameLength(saveStream, 'PLGI', 4)
        self.plgi2 = NameLength(saveStream, 'PLGI', 8)
        self.date2 = NameLength(saveStream, 'DATE', 84)
        self.date3 = NameLength(saveStream, 'DATE', 84)
        # There seem to be 8 bytes here; guessing two integers
        (self.integer1,) = struct.unpack('i', saveStream.read(4))
        (self.integer2,) = struct.unpack('i', saveStream.read(4))
        self.cnsl = NameLength(saveStream, 'CNSL', 228)
        self.wrld1 = NameLength(saveStream, 'WRLD', 2)
        (num_continents,) = struct.unpack_from('h', self.wrld1.data)
        self.wrld2 = NameLength(saveStream, 'WRLD', 164)
        self.mapHeight = struct.unpack_from('41i', self.wrld2.data)[1]
        self.mapWidth = struct.unpack_from('41i', self.wrld2.data)[6]
        #print "map: " + str(self.mapWidth) + " x " + str(self.mapHeight)
        self.wrld2 = NameLength(saveStream, 'WRLD', 52)
        self.tiles = []
        for tile in range(self.mapWidth / 2 * self.mapHeight):
            data = []
            for i in range(4):
                data.append(NameLength(saveStream, 'TILE'))
            self.tiles.append(data)
        self.continents = []
        for i in range(num_continents):
            self.continents.append(NameLength(saveStream, 'CONT'))
        print self.continents[-1]

        # There is some data--usually of length 0x68--here that looks like an integer array to me
        # I believe the number of integers is the number of resources including bonus resources
        self.integers = []
        data  = ""
        data = saveStream.read(4)
        while data <> 'LEAD':
            (integer,) = struct.unpack_from('i', data)
            self.integers.append(integer)
            data = saveStream.read(4)
        #print len(self.integers)
        #print self.integers

        # this isn't right
        #self.lead1 = NameLength(saveStream, 'LEAD')
        #print self.lead1

        self.leads = []
        for i in range(32):
            data = []
            # SKIPPING over LEAD section since I haven't figured it out yet
            self.leadLength = horspool.boyermoore_horspool(saveStream, "CULT")
            print 'LEAD section length: {0}'.format(self.leadLength)

            data.append(HorspoolNameLength(saveStream, 'CULT', 0x10))
            data.append(NameLength(saveStream, 'ESPN', 0x20))
            data.append(NameLength(saveStream, 'ESPN', 0x20))
#            for j in range(len(data)):
#                print data[j]
            self.leads.append(data)

        self.units = []
        # SKIPPING to first UNIT because I'm not sure what's here
        self.skipLength = horspool.boyermoore_horspool(saveStream, "UNIT")
        print 'pre-UNIT skipped length: {0}'.format(self.leadLength)
        unit = []
        unit.append(HorspoolNameLength(saveStream, 'UNIT', 0x1d8))
        # TODO : read IDLS and append them to UNIT array
        idls = Idls(saveStream)
        print idls.records
        unit.append(idls)
        self.units.append(unit)

        nextSection = NameLength(saveStream, 'UNIT', 0x1d8)
        print nextSection.name
        while nextSection.name == 'UNIT':
            unit = []
            unit.append(nextSection)
            idls = Idls(saveStream)
            #print idls.records
            unit.append(idls)
            self.units.append(unit)
            nextSection = NameLength(saveStream, 'UNIT', 0x1d8)
        #print nextSection

        # TODO: what about saves with no cities? (4000 BC)
        self.cities = []
    #while nextSection.name == 'CITY':
        city = []
        city.append(nextSection)
        city.append(NameLength(saveStream, 'CITY'))
        city.append(NameLength(saveStream, 'CITY'))
        city.append(NameLength(saveStream, 'CITY'))
        city.append(NameLength(saveStream, 'CITY'))
        city.append(Popd(saveStream, 'POPD'))
        binfLength = horspool.boyermoore_horspool(saveStream, "BITM")
        print "BINF Length:", binfLength
        city.append(HorspoolNameLength(saveStream, 'BITM'))
        city.append(NameLength(saveStream, 'DATE'))

        self.cities.append(city)
        #nextSection = NameLength(saveStream, 'CITY')
    #print nextSection

        ## TODO: There are a lot of empty CITY sections. Need to figure out format
        ## Also, some cities have CTPG sections I'm not reading in
        ## Also, CLNY sections, if any, would be between CITY and PALV

        skipToPalvLength = horspool.boyermoore_horspool(saveStream, "PALV")
        self.palv = []
        self.palv.append(HorspoolNameLength(saveStream, 'PALV'))
        #print self.palv[-1]
        for i in range(31):
            self.palv.append(NameLength(saveStream, 'PALV'))
            #print self.palv[-1]

        # HIST
        # Already there after PALV, but don't see how it's formatted

        # TUTR
        print "Skip to TUTR: ", horspool.boyermoore_horspool(saveStream, "TUTR")
        self.tutr = HorspoolNameLength(saveStream, 'TUTR')
        #print self.tutr

        # FAXX
        self.faxx = NameLength(saveStream, 'FAXX')
        #print self.faxx

        # RPLS
        (rplsname, numRecords) = struct.unpack_from('4si', saveStream.read(8))
        print rplsname, numRecords
        # this isn't right
        print NameLength(saveStream)
#        self.rpls = []
#        #for i in range(numRecords)
#        for i in range(5):
#            self.rpls.append(NameLength(saveStream))
#            print self.rpls[-1]

        print "\nWhat's Next:"
        self.whatsnext = hexdump(saveStream.read(0x800))
        print self.whatsnext

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
