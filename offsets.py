#!/usr/bin/env python
# -*- coding: latin-1 -*-

import sys
from horspool import horspool    # to seek to first match; from http://inspirated.com/2010/06/19/using-boyer-moore-horspool-algorithm-on-file-streams-in-python

def seek(saveStream, string):
      reloffset = horspool.boyermoore_horspool(saveStream, string)
      print string + " relative offset: " + str(reloffset)
      return reloffset

def main():
    """Analyzing relative offsets between data structures"""
    saveStream = sys.stdin
    reloffset = 0
    string="LEAD"
    reloffset = seek(saveStream, "GAME")
    reloffset = seek(saveStream, "GAME")
    reloffset = seek(saveStream, "DATE")
    reloffset = seek(saveStream, "PLGI")
    reloffset = seek(saveStream, "PLGI")
    reloffset = seek(saveStream, "DATE")
    reloffset = seek(saveStream, "DATE")
    reloffset = seek(saveStream, "CNSL")
    reloffset = seek(saveStream, "WRLD")
#    while reloffset >= 0 :
#      #reloffset = horspool.boyermoore_horspool(saveStream, string)
#      #print string + " relative offset: " + str(reloffset)
#      reloffset = seek(saveStream, "LEAD")
#      reloffset = seek(saveStream, "CULT")
#      reloffset = seek(saveStream, "CULT")
    #saveStream.close()
    return

if __name__=="__main__":
    main()
