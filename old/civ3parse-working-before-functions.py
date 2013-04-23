#!/usr/bin/env python

# 2013-04-19 Start of attempt to parse a decompressed Civ3 save file
# unc-test.sav

import struct		# For parsing runs of binary data

print "Hello, World"

message = "This is the message"

print message

saveFilePath = "unc-test.sav"

saveFile = open(saveFilePath, 'rb')

print "Skipping first 30 byts because the CIV3 header doesn't seem to be followed by a length"
saveFile.seek(30,0)

# way to loop until EOF adapted from http://stackoverflow.com/questions/1752107/how-to-loop-until-eof-in-python
for buffer in iter(lambda: saveFile.read(8), ''):
	(section,) = struct.unpack('4s',buffer[0:4])
	(sectionLength,) = struct.unpack('i',buffer[4:8])
	print "Section Name:", section, ", Length: ", sectionLength
	saveFile.seek(sectionLength,1)
