#!/usr/bin/env python

import json

blah = []

for y in range(10):
    blah.append([])
    for x in range(10):
        blah[y].append({'this':'that ' + str(x) + ',' + str(y)})

print blah
print str(blah)
print "and now the JSON..."
print json.dumps(blah)
print json.JSONEncoder(skipkeys=True).encode(blah)
