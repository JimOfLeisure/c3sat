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

# 2013-04-19 Start of attempt to parse a decompressed Civ3 save file
# unc-test.sav

import struct		# For parsing runs of binary data
import sys		# for exiting program during development/debugging
import inspect		# for listing calling function name in debug prints

def readLength(saveStream):	# I repeatedly need to read a 4-byte length integer from the file stream
	buffer = saveStream.read(4)
	(length,) = struct.unpack('i',buffer[0:4])
	return length

def printDebug(saveStream, debugInfo):
	readPosition = saveStream.tell()
	print 'Debug ' \
	+ inspect.stack()[1][3] \
	+ ' ' \
	+ debugInfo \
	+ ' hex offset ' \
	+ hex(readPosition) \
	+ ' ' \
	+ ' ' \

def skipBytes(saveStream, skipNum):
	printDebug(saveStream, 'Skipping ' + hex(skipNum) + ' bytes')
	saveStream.seek(skipNum,1)

def skipToOffset(saveStream, skipTo):
	printDebug(saveStream, 'Skipping to offset ' + hex(skipTo))
	saveStream.seek(skipTo,0)

def sectionGeneric(saveStream, sectionName):	# Reads length, then skips that many bytes
	sectionLength = readLength(saveStream)
#	print "sectionGeneric says:" , "Section Name:", sectionName, ", Length: ", sectionLength
	printDebug(saveStream, sectionName + ' length ' + str(sectionLength))
	saveStream.seek(sectionLength,1)

def sectionSubRecord(saveStream, sectionName):		# Reads length in subrecords then iterates through them
	numSubRecords = readLength(saveStream)
	print "sectionSubRecord says:", "Section Name:", sectionName, "Number of subrecords:", numSubRecords
	printDebug(saveStream, sectionName + ' numSubRecords ' + str(numSubRecords))
#	sys.exit("Stopping point in development")
	while numSubRecords > 0:
		sectionGeneric(saveStream, sectionName)
		numSubRecords -= 1


# Replacing with generic skip function
#def sectionCIV3(saveStream):	# I don't know much about this section yet
#	print "Skipping 26 bytes because the CIV3 header doesn't seem to be followed by a length"
#	saveStream.seek(26,1)

def sectionBICQ(saveStream):	# This section needs slightly different handling; read the "VER#" then iterate subrecords
	buffer = saveStream.read(4)
	(thisSaysVerNum,) = struct.unpack('4s', buffer[0:4])
	print "sectionBICQ says: This should say VER#:", thisSaysVerNum
	sectionSubRecord(saveStream, "BICQVER#")

# using a skip function here temporarily
def sectionGAME(saveStream):	# A more complex system of subrecords
#	numScenarioProperties= readLength(saveStream)
#	while numScenarioProperties > 0:
#		scenarioPropertyLength =  readLength(saveStream)
#		print "Scenario Property #", numScenarioProperties, "Length:",scenarioPropertyLength
#		numScenarioProperties -= 1
#	sys.exit("Stopping point in development")
	print "Skipping 11384 bytes because I can\'t figure out the GAME section"
	saveStream.seek(11384,1)

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

def parseSave():
	saveFilePath = "unc-test.sav"
	
	saveFile = open(saveFilePath, 'rb')

	gameSectionsRead = 0

	# way to loop until EOF adapted from http://stackoverflow.com/questions/1752107/how-to-loop-until-eof-in-python
	for buffer in iter(lambda: saveFile.read(4), ''):
		(sectionName,) = struct.unpack('4s',buffer[0:4])
		#printDebug(saveFile, 'just read name')
		# Ensure sectionName is ASCII, taken from http://stackoverflow.com/questions/196345/how-to-check-if-a-string-in-python-is-in-ascii
		try:
			sectionName.decode('ascii')
		except UnicodeDecodeError:
			readPosition = saveFile.tell()
			errorMessage = 'ERROR: parseSave(): sectionName is not an ASCII string. Offset:' + str(readPosition) + ' Hex: ' + hex(readPosition)
			sys.exit(errorMessage)
		# Apparently I need to learn polymorphism but for now I can make it work with chained if-then-else
		if sectionName == 'CIV3':
			#sectionCIV3(saveFile)
			skipBytes(saveFile, 26)
		elif sectionName == 'BICQ':
			sectionBICQ(saveFile)
		elif sectionName == 'GAME':
			#sectionGAME(saveFile)
			if gameSectionsRead > 0:
				sectionGeneric(saveFile, sectionName)
			else:
				skipToOffset(saveFile, 0x32c6)	# Skipping; assuming this offset will not work with arbitrary save files
			gameSectionsRead += 1
		else:
			sectionGeneric(saveFile, sectionName)

def dumpWholeFile():
	saveFilePath = "unc-test.sav"
	saveFile = open(saveFilePath, 'rb')
	#hexdump(saveFile)
	# way to loop until EOF adapted from http://stackoverflow.com/questions/1752107/how-to-loop-until-eof-in-python
	for buffer in iter(lambda: saveFile.read(4), ''):
		print hexdump(buffer,4)


def main():
	#parseSave()
    dumpWholeFile()

main()
