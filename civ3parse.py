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

# 2013-04-19 Start of attempt to parse a decompressed Civ3 save file
# unc-test.sav

import re		# Regular Expression comparison
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


def parseSave():
	saveFilePath = "gamesaves/unc-test.sav"
	
	saveFile = open(saveFilePath, 'rb')

	gameSectionsRead = 0

	# way to loop until EOF adapted from http://stackoverflow.com/questions/1752107/how-to-loop-until-eof-in-python
	for buffer in iter(lambda: saveFile.read(4), ''):
		(sectionName,) = struct.unpack('4s',buffer[0:4])
		match = re.match('[A-Z, ]',sectionName)
		if match is None:
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

def main():
	parseSave()

main()
