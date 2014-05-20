#!/usr/bin/env python

import locale
import os
import sys
import urllib2

def boyermoore_horspool(fd, needle):
    nlen = len(needle)
    nlast = nlen - 1

    skip = []
    for k in range(256):
        skip.append(nlen)
    for k in range(nlast):
        skip[ord(needle[k])] = nlast - k
    skip = tuple(skip)

    pos = 0
    consumed = 0
    haystack = bytes()
    while True:
        more = nlen - (consumed - pos)
        morebytes = fd.read(more)
        haystack = haystack[more:] + morebytes

        if len(morebytes) < more:
            return -1
        consumed = consumed + more

        i = nlast
        while i >= 0 and haystack[i] == needle[i]:
            i = i - 1
        if i == -1:
            return pos

        pos = pos + skip[ord(haystack[nlast])]

    return -1

if __name__ == "__main__":
    if len(sys.argv) < 3:
        print "Usage: horspool.py <url> <search text>"
        sys.exit(-1)

    url = sys.argv[1]
    needle = sys.argv[2]
    needle = needle.decode('string_escape')

    fd = urllib2.urlopen(url)
    offset = boyermoore_horspool(fd, needle)
    print hex(offset), '::', offset
    fd.close()

