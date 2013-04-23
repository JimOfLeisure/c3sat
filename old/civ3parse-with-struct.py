#!/usr/bin/env python

# 2013-04-19 Jim Nelson . Start of attempt to parse a decompressed Civ3 save file
# unc-test.sav

import struct		# For parsing runs of binary data

print "Hello, World"

message = "This is the message"

print message

saveFilePath = "unc-test.sav"

saveFile = open(saveFilePath, 'rb')

buffer = saveFile.read(4)
section = struct.unpack('4s',buffer)

print section
