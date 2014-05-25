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

import wrld
import urllib2
import subprocess
import os

def application(env, start_response):
    """This module instantiates wrld.parse_save() and returns an svg file for the map"""
    #start_response('200 OK', [('Content-Type','text/html')])
    #return "Hello World From Python"
    testInputFile = "gamesaves/unc-test.sav"
    #saveFile = open(testInputFile, 'rb')

    #testInputUrl = "http://lib.bigmoneyjim.com/civfan/Mao%20of%20the%20Chinese,%20130%20AD.SAV"
    #testInputUrl = "http://lib.bigmoneyjim.com/civfan/Puppeteer-joinworker-Mongols-Reroll-as-Regent-4000%20BC.SAV"
    #testInputUrl = "file:gamesaves/unc-test.sav"
    #testInputUrl = "file:Saves/end turn conquest English, 1340 AD.SAV"
    testInputUrl = "file:Puppeteer-joinworker-Mongols-Reroll-as-Regent-4000 BC.SAV"

    # If first 4 bytes aren't "CIV3" then we should decompress
#    saveFile = urllib2.urlopen(testInputUrl)
#    buffer = saveFile.read(4)
#    saveFile.close()
    #if buffer == 'CIV3':
    if False:
        saveFile = urllib2.urlopen(testInputUrl)
        game = wrld.parse_save(saveFile)
    else:
        print "Stub. Put decompress code here"
        myCompressedFile = urllib2.urlopen(testInputUrl)
        print myCompressedFile.getcode()
        print myCompressedFile.geturl()
        print myCompressedFile.info()
        #process = subprocess.Popen(['./blast'], stdin=myCompressedFile, stdout=subprocess.PIPE,stderr=subprocess.PIPE, shell=True)
        process = subprocess.Popen(['./blast'], stdin=myCompressedFile, stdout=subprocess.PIPE,stderr=subprocess.PIPE)
        #process = subprocess.Popen(['./blast'], stdin=subprocess.PIPE, stdout=subprocess.PIPE,stderr=subprocess.PIPE)
        #process.stdin.write(myCompressedFile.read())
        #process.stdin.write(myCompressedFile)
        print "Hi"
        #process.stdin.close()
        #myCompressedFile.close()
        #buffer = process.stdout.read(4)
        #print buffer
#        out = process.communicate(myCompressedFile.read())
        #out = process.communicate()
        #print len(out[0])
#        print out[0][0]
#        print out[0][1]
#        print out[0][2]
#        print out[0][3]
        #print out[1]
        #saveFile = process.stdout
        game = wrld.parse_save(process.stdout)
        #game = wrld.parse_save(process.communicate()[0])



    start_response('200 OK', [('Content-Type','image/svg+xml')])
    return game.Tiles.svg_out(True)

    saveFile.close()
