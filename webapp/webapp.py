#!/usr/bin/env python

# Copyright (c) 2014 Jim Nelson
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

import readciv3
import urllib
import subprocess
import os
import urlparse

def application(env, start_response):
    """This module instantiates readciv3.parse_save() and returns an svg file for the map"""
    #start_response('200 OK', [('Content-Type','text/html')])
    #return "Hello World From Python"

    #InputUrl = "http://lib.bigmoneyjim.com/civfan/Mao%20of%20the%20Chinese,%20130%20AD.SAV"
    #InputUrl = "http://lib.bigmoneyjim.com/civfan/Puppeteer-joinworker-Mongols-Reroll-as-Regent-4000%20BC.SAV"
    #InputUrl = "http://lib.bigmoneyjim.com/civfan/c3sat/unc-test.sav"   # uncompressed save
    #InputUrl = "http://lib.bigmoneyjim.com/civfan/" # intentional non-save url

    myargs =  urlparse.parse_qs(env["QUERY_STRING"])
    #print myargs
    InputUrl = myargs["url"][0]
    #print InputUrl

    # If first 4 bytes aren't "CIV3" then we should decompress
    saveFile = urllib.urlopen(InputUrl)
    buffer = saveFile.read(4)
    saveFile.close()
    #print buffer

    if buffer == 'CIV3':
        #print "Uncompressed file"
        saveFile = urllib.urlopen(InputUrl)
    else:
        #print "Decompress"
        myCompressedFile = urllib.urlopen(InputUrl)
        #print myCompressedFile.getcode()
        #print myCompressedFile.geturl()
        #print myCompressedFile.info()

        process = subprocess.Popen(['./blast'], stdin=myCompressedFile, stdout=subprocess.PIPE,stderr=subprocess.PIPE)

        #print process.stderr.read()
        #print process.stdout.read()
        saveFile = process.stdout
        
#    for k in env:
#        print k, env[k]

    game = readciv3.parse_save(saveFile)
    saveFile.close()
    start_response('200 OK', [('Content-Type','image/svg+xml')])
    return game.Tiles.svg_out()

