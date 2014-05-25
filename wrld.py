#!/usr/bin/env python
# -*- coding: latin-1 -*-

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

# 2014-05-19 Now I want to skip to the first WRLD section in arbitrary uncompressed SAV files

import struct	# For parsing binary data
import json     # to export JSON for the HTML browser
import horspool    # to seek to first match; from http://inspirated.com/2010/06/19/using-boyer-moore-horspool-algorithm-on-file-streams-in-python
import sys

class GenericSection:
    """Base class for reading SAV sections."""
    def __init__(self, saveStream):
        buffer = saveStream.read(8)
        (self.name, self.length,) = struct.unpack_from('4si', buffer)
        #self.offset = saveStream.tell()
        self.buffer = saveStream.read(self.length)

class Tile:
    """Class for each logical tile."""
    def __init__(self, saveStream):

        self.info = {}

        self.Tile36 = GenericSection(saveStream)
        self.continent = get_short(self.Tile36.buffer, 0x1a)
        self.info['continent'] = get_short(self.Tile36.buffer, 0x1a)
        del self.Tile36.buffer

        self.Tile12 = GenericSection(saveStream)
        self.info['terrain'] = get_byte(self.Tile12.buffer, 0x5)
        #self.whatsthis = get_byte(self.Tile12.buffer, 0xa)
        self.whatsthis = get_byte(self.Tile12.buffer, 0x5)
        self.whatsthis2 = get_byte(self.Tile12.buffer, 0x5)
        self.whatsthis3 = get_byte(self.Tile12.buffer, 0xa)
        self.whatsthis4 = get_byte(self.Tile12.buffer, 0xb)
        self.whatsthis5 = get_int(self.Tile12.buffer, 0x0)
        del self.Tile12.buffer

        self.Tile4 = GenericSection(saveStream)
        del self.Tile4.buffer

        self.Tile128 = GenericSection(saveStream)
        self.is_visible_to = get_int(self.Tile128.buffer, 0)
        self.is_visible_now_to = get_int(self.Tile128.buffer, 4)
        self.is_visible = self.is_visible_to & 0x02
        #self.is_visible = self.is_visible_to & 0x10
        self.is_visible_now = self.is_visible_now_to & 0x02
        del self.Tile128.buffer

class Tiles:
    """Class to read all tiles"""
    def __init__(self, saveStream, width, height):
        self.width = width      # These may eventually be redundant to a parent class
        self.height = height
        self.tile = []          # List of individual tiles
        self.tile_matrix = []       # x,y matrix of individual tiles
        self.tile_iso_matrix = []   # faux isometric padded matrix  of individual tiles
#        logical_tiles = width / 2 * height
#        while logical_tiles > 0:
#            self.tile.append(Tile(saveStream))
#            logical_tiles -= 1
        for y in range(height):
            self.tile_matrix.append([])
            self.tile_iso_matrix.append([])
            if y % 2 == 1:
                self.tile_iso_matrix[y].append(None)
            for x in range(width / 2):
                this_tile = Tile(saveStream)

                self.tile.append(this_tile)

                self.tile_matrix[y].append(this_tile)

                self.tile_iso_matrix[y].append(this_tile)
                self.tile_iso_matrix[y].append(None)

            if y % 2 == 0:
                self.tile_iso_matrix[y].append(None)

    def map_id(self, i):
        """Return a string to be used as a CSS ID for the tile group. i is the index of self.tile"""
        return 'map' + str(i)

    def svg_out(self, spoiler=False):
        """Return a string of svg-coded map"""
        x_axis_wrap = True
        y_axis_wrap = False
        tile_width = 128
        tile_height = 64
        map_width = (self.width * tile_width / 2) + (tile_width / 2)
        map_height = (self.height * tile_height / 2) + (tile_height / 2)
        svg_string = ""
        svg_string += '<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" version="1.1" x="0" y="0" viewBox="0 0 ' + str(map_width) + ' ' + str(map_height) + '">\n'
        #svg_string += '<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" version="1.1" width="' + str(map_width) + '" height="' + str(map_height) + '">\n'
        svg_string += """
<defs>
  <g transform="scale(0.1) translate(-355,-450)" id="myTree">
    <rect
       style="fill:#502d16"
       id="rect3755"
       width="38.385796"
       height="135.36044"
       x="333.35034"
       y="412.93561" />
    <path
       style="fill:#003e00;stroke:#000000;stroke-width:1px;stroke-linecap:butt;stroke-linejoin:miter;stroke-opacity:1"
       d="m 350.52294,162.41779 c 45.35523,49.58592 68.30881,76.28117 137.38074,106.06602 -49.49747,3.16986 -65.65992,7.45477 -115.15739,-8.08122 37.19346,33.65589 79.36925,63.67564 136.3706,89.90357 -29.52448,8.86132 -77.90225,-10.22203 -126.26907,-20.20305 31.03325,29.4296 63.72488,75.96279 115.15739,105.05586 -84.46056,-24.01069 -184.4487,-22.38115 -314.15745,10e-6 40.71838,-21.12578 104.53494,-63.53265 141.42137,-105.05587 -33.89684,20.76304 -87.0985,20.089 -133.34014,16.16243 54.92101,-23.31175 90.86412,-44.06567 135.36045,-83.84266 -20.68253,3.88388 -77.82911,7.71658 -107.07619,10.10154 26.31048,-5.29696 125.20682,-84.59227 130.30969,-110.10663 z" />
  </g>    
  <g id="myForest">    
    <use xlink:href = "#myTree" x="-10" y="-2" />
    <use xlink:href = "#myTree" x="10" y="2" />
    <use xlink:href = "#myTree" x="-40" y="7" />
    <use xlink:href = "#myTree" x="32" y="7" />
  </g>    
    <linearGradient
       id="linearGradient5227">
      <stop
         id="stop5229"
         style="stop-color:#008000;stop-opacity:1"
         offset="0" />
      <stop
         id="stop5235"
         style="stop-color:#006700;stop-opacity:1"
         offset="0.84402692" />
      <stop
         id="stop5231"
         style="stop-color:#008000;stop-opacity:1"
         offset="1" />
    </linearGradient>
    <linearGradient
       id="linearGradient5222">
      <stop
         id="stop5224"
         style="stop-color:#000000;stop-opacity:1"
         offset="0" />
      <stop
         id="stop5226"
         style="stop-color:#008000;stop-opacity:1"
         offset="1" />
    </linearGradient>
    <linearGradient
       id="linearGradient5212">
      <stop
         id="stop5214"
         style="stop-color:#000000;stop-opacity:1"
         offset="0" />
    </linearGradient>
    <radialGradient
       cx="76.143784"
       cy="962.11768"
       r="60.366074"
       fx="76.143784"
       fy="962.11768"
       id="radialGradient5233"
       xlink:href="#linearGradient5227"
       gradientUnits="userSpaceOnUse"
       gradientTransform="matrix(-0.98699431,0.64563536,-0.46697231,-0.71386896,600.28401,1624.7186)" />
    <linearGradient
       x1="64.732147"
       y1="985.2901"
       x2="64.107155"
       y2="1018.683"
       id="linearGradient5247"
       xlink:href="#linearGradient5222"
       gradientUnits="userSpaceOnUse" />

  <path
     d="m 10.446429,1019.0586 c 4.044564,-11.572 19.854103,-23.29369 38.482142,-30.35713 14.155388,-5.36749 22.120429,-2.86384 30.178573,-0.17856 10.868126,3.62166 25.732176,13.89989 36.964296,32.67849 -16.238763,11.6432 -32.832015,22.8648 -51.250011,21.6072 -31.562398,-2.1549 -44.640301,-15.9149 -54.375,-23.75 z"
     id="myHill"
     transform="translate(-64,-1020.36218)"
     style="fill:url(#radialGradient5233);fill-opacity:1;stroke:url(#linearGradient5247)" />
  <path
     d="M 8.0812204,31.603036 C 13.324636,11.111392 18.747476,3.5155048 30.557114,-16.37921 38.411455,-5.9814302 37.263782,-7.4855782 41.92133,2.2728428 52.678471,-10.824571 51.489681,-14.65148 57.831233,-34.056879 57.039687,-34.848425 68.6306,-19.150733 71.468292,-18.399515 74.682439,-17.548638 90.655911,-37.909406 89.903576,-38.09749 96.851004,-2.0203051 103.6441,0 116.42008,31.09796 c -10.80604,5.157699 -20.7117,11.108946 -36.112953,13.889597 -9.313114,1.681455 -18.393434,2.222457 -33.082496,0 C 32.44483,42.751372 20.311577,37.290714 8.0812204,31.603036 z"
     id="myMountain"
     transform="translate(-64,-32)"
     style="fill:#803300;stroke:#000000;stroke-width:1px" />
    <linearGradient
       id="linearGradient4298">
      <stop
         id="stop4300"
         style="stop-color:#000000;stop-opacity:1"
         offset="0" />
      <stop
         id="stop4302"
         style="stop-color:#803300;stop-opacity:1"
         offset="1" />
    </linearGradient>
    <linearGradient
       id="linearGradient4288">
      <stop
         id="stop4290"
         style="stop-color:#000000;stop-opacity:1"
         offset="0" />
      <stop
         id="stop4292"
         style="stop-color:#803300;stop-opacity:1"
         offset="1" />
    </linearGradient>
    <radialGradient
       cx="59.122124"
       cy="-28.951275"
       r="54.66943"
       fx="59.122124"
       fy="-28.951275"
       id="radialGradient4294"
       xlink:href="#linearGradient4288"
       gradientUnits="userSpaceOnUse"
       gradientTransform="matrix(-0.75311906,1.0400717,-0.84152356,-0.60934973,84.43377,-111.3177)" />
    <radialGradient
       cx="69.945213"
       cy="-45.220779"
       r="23.48097"
       fx="69.945213"
       fy="-45.220779"
       id="radialGradient4304"
       xlink:href="#linearGradient4298"
       gradientUnits="userSpaceOnUse"
       gradientTransform="matrix(-0.87568994,1.0988247,-1.1618204,-0.92589342,85.736026,-154.75086)" />
<g id="myVolcano" transform="translate(-64,-32)">
  <path
     d="m 44.547527,-33.20506 c 8.275665,-5.039021 10.202241,-7.763882 19.68397,-8.625252 8.990378,-0.841691 16.716875,0.02087 25.117641,3.049646 L 67.427682,-15.36905 z"
     id="path4296"
     style="fill:url(#radialGradient4304);fill-opacity:1;stroke:#000000;stroke-width:1px;stroke-linecap:butt;stroke-linejoin:miter;stroke-opacity:1" />
  <path
     d="M 8.0812204,31.603036 C 19.380724,0 32.777118,0 43.941636,-32.54165 c 2.199525,1.105775 11.580429,9.211579 18.940359,9.848987 7.111092,0.615856 27.021581,-15.404827 27.021581,-15.404827 C 96.851004,-2.0203051 103.6441,0 116.42008,31.09796 c -10.80604,5.157699 -20.7117,11.108946 -36.112953,13.889597 -9.313114,1.681455 -18.393434,2.222457 -33.082496,0 C 32.44483,42.751372 20.311577,37.290714 8.0812204,31.603036 z"
     id="path2984"
     style="fill:url(#radialGradient4294);fill-opacity:1;stroke:#000000;stroke-width:1px;stroke-linecap:butt;stroke-linejoin:miter;stroke-opacity:1" />

</g>
<polygon id="isometbrick" points="0,16 128,16 128,48 0,48" />
<polygon id="isotile" points="0,32 64,0 128,32 64,64" />
<use id="maptileshape" xlink:href = "#isotile" />
<use id="desert" xlink:href="#maptileshape" style="fill:cornsilk" />
<use id="plains" xlink:href="#maptileshape" style="fill:sandybrown" />
<use id="grassland" xlink:href="#maptileshape" style="fill:green" />
<use id="tundra" xlink:href="#maptileshape" style="fill:white" />
<use id="coast" xlink:href="#maptileshape" style="fill:blue" />
<use id="sea" xlink:href="#maptileshape" style="fill:mediumblue" />
<use id="ocean" xlink:href="#maptileshape" style="fill:darkblue" />
<use id="fog" xlink:href="#maptileshape" style="fill:black" />
  </defs>
"""
        svg_string += '<rect class="mapEdge" x="0" y="0" width="' + str(map_width) + '" height="' + str(map_height) + '" />\n'
        for y in range(self.height):
            x_indent = (y % 2) * tile_width / 2
            y_offset = y * tile_height / 2
            svg_string += '<g row="' + str(y) + '" transform="translate(' + str(x_indent) + ', ' + str(y_offset) + ')">\n'
            for x in range(self.width / 2):
                i = x  + y * self.width /2
                info = hex(self.tile[i].whatsthis)
                cssclass = 'tile'
                if 0 <= i < len(self.tile):
                    svg_string += '  <g id="' + self.map_id(i) + '" transform="translate(' + str(x * tile_width) + ', 0)">\n'
                    if self.tile[i].is_visible or spoiler:
                        # Get right-nibble of terrain byte
                        base_terrain = self.tile[i].info['terrain'] & 0x0F
                        # Get left-nibble of terrain byte: bit-rotate right 4, then mask to be sure it wasn't more than a byte
                        overlay_terrain = (self.tile[i].info['terrain'] >> 4) & 0x0F
                        cssclass += ' baseterrain'
                        if base_terrain == 0:
                            svg_string += '<use xlink:href="#desert" class="' + cssclass +'" />\n'
                        elif base_terrain == 1:
                            svg_string += '<use xlink:href="#plains" class="' + cssclass +'" />\n'
                        elif base_terrain == 2:
                            svg_string += '<use xlink:href="#grassland" class="' + cssclass +'" />\n'
                        elif base_terrain == 3:
                            svg_string += '<use xlink:href="#tundra" class="' + cssclass +'" />\n'
                        elif base_terrain == 11:
                            svg_string += '<use xlink:href="#coast" class="' + cssclass +'" />\n'
                        elif base_terrain == 12:
                            svg_string += '<use xlink:href="#sea" class="' + cssclass +'" />\n'
                        elif base_terrain == 13:
                            svg_string += '<use xlink:href="#ocean" class="' + cssclass +'" />\n'

                        # Not sure I need this comparison; may just be able to key off the nibble value
                        if overlay_terrain <> base_terrain:
                            cssclass = 'overlayterrain terroverlay' + str(overlay_terrain)
                            if overlay_terrain == 0x04:
                                # Flood plain
                                svg_string += '    <text class="' + cssclass + '" text-anchor="middle" alignment-baseline="central" x="' + str(tile_width / 2) + '" y="' + str(tile_height / 2) + '">Flood Plain</text>\n'
                            elif overlay_terrain == 0x05:
                                # Hill
                                #svg_string += '    <text class="' + cssclass + '" text-anchor="middle" alignment-baseline="central" x="' + str(tile_width / 2) + '" y="' + str(tile_height / 2) + '">○</text>\n'
                                svg_string += '    <use class="' + cssclass + '" x="' + str(tile_width / 2) + '" y="' + str(tile_height / 2) + '" xlink:href = "#myHill" />\n'
                            elif overlay_terrain == 0x06:
                                # Mountain
                                #svg_string += '    <text class="' + cssclass + '" text-anchor="middle" alignment-baseline="central" x="' + str(tile_width / 2) + '" y="' + str(tile_height / 2) + '">▲</text>\n'
                                svg_string += '    <use class="' + cssclass + '" x="' + str(tile_width / 2) + '" y="' + str(tile_height / 2) + '" xlink:href = "#myMountain" />\n'
                            elif overlay_terrain == 0x07:
                                # Forest
                                #svg_string += '    <text class="' + cssclass + '" text-anchor="middle" alignment-baseline="central" x="' + str(tile_width / 2) + '" y="' + str(tile_height / 2) + '">↑↑↑</text>\n'
                                #svg_string += '    <use class="' + cssclass + '" x="' + str(tile_width / 2) + '" y="' + str(tile_height / 2) + '" xlink:href = "#myTree" />\n'
                                svg_string += '    <use class="' + cssclass + '" x="' + str(tile_width / 2) + '" y="' + str(tile_height / 2) + '" xlink:href = "#myForest" />\n'
                            elif overlay_terrain == 0x08:
                                # Jungle
                                svg_string += '    <text class="' + cssclass + '" text-anchor="middle" alignment-baseline="central" x="' + str(tile_width / 2) + '" y="' + str(tile_height / 2) + '">Jungle</text>\n'
                            elif overlay_terrain == 0x09:
                                # Marsh
                                svg_string += '    <text class="' + cssclass + '" text-anchor="middle" alignment-baseline="central" x="' + str(tile_width / 2) + '" y="' + str(tile_height / 2) + '">Marsh</text>\n'
                            elif overlay_terrain == 0x0a:
                                # Volcano
                                #svg_string += '    <text class="' + cssclass + '" text-anchor="middle" alignment-baseline="central" x="' + str(tile_width / 2) + '" y="' + str(tile_height / 2) + '">Volcano</text>\n'
                                svg_string += '    <use class="' + cssclass + '" x="' + str(tile_width / 2) + '" y="' + str(tile_height / 2) + '" xlink:href = "#myVolcano" />\n'
                            else:
                                svg_string += '    <text class="' + cssclass + '" text-anchor="middle" alignment-baseline="central" x="' + str(tile_width / 2) + '" y="' + str(tile_height / 2) + '">' + info + '</text>\n'
                        svg_string += '    <text class="whatsthis" style="display:none" text-anchor="middle" alignment-baseline="central" x="' + str(tile_width / 2) + '" y="' + str(tile_height / 2) + '">' + info + '</text>\n'
                        #svg_string += '  </g>\n'
                    else:
                        cssclass = 'fog'
                        svg_string += '<use xlink:href="#fog" class="' + cssclass +'" />\n'
                        #svg_string += '  </g>\n'
                    svg_string += '  </g>\n'
            # link the first item and place at the end for even rows; link to the last item and place at the first. Will be half-cropped by viewport
            # using math (even lines have 0 remainder, multiplying to cancel out values) instead of if, but it's a little harder to follow
            svg_string += '  <use xlink:href="#' + self.map_id((y * self.width / 2) + (x * (y % 2))) + '" transform="translate(' + str((map_width - tile_width / 2) - (map_width - tile_width / 2) * 2 * (y % 2)) + ', 0)" />\n'
            svg_string += '</g>\n'
        svg_string += '</svg>\n'
        return svg_string

class Wrld:
    """Class for 3 WRLD sections"""
    def __init__(self, saveStream):
        """Currently calling this from the horspool seek, so WRLD is already consumed from the stream. Read the length first."""
        self.name = "WRLD"
        buffer = saveStream.read(4)
        (self.length,) = struct.unpack_from('i', buffer)
        #print self.length
        self.buffer = saveStream.read(self.length)
        # Extract any data here, but I think it's only 2 bytes
        #print self.name
        #print hexdump(self.buffer)
        del self.buffer

        self.Wrld2 = GenericSection(saveStream)
        #print self.Wrld2.name
        self.values = struct.unpack_from('41i', self.Wrld2.buffer)
        self.height = self.values[1]
        self.width = self.values[6]
        #print self.height
        #print self.width
        #print self.values
        #print hexdump(self.Wrld2.buffer)
        del self.Wrld2.buffer

        self.Wrld3 = GenericSection(saveStream)
        #print self.Wrld3.name
        #print hexdump(self.Wrld3.buffer)
        del self.Wrld3.buffer

        self.Tiles = Tiles(saveStream, self.width, self.height)


def get_byte(buffer, offset):
    """Unpack an byte from a buffer at the given offest."""
    (the_byte,) = struct.unpack('B', buffer[offset:offset+1])
    return the_byte

def get_short(buffer, offset):
    """Unpack an int from a buffer at the given offest."""
    (the_short,) = struct.unpack('H', buffer[offset:offset+2])
    return the_short

def get_int(buffer, offset):
    """Unpack an int from a buffer at the given offest."""
    (the_int,) = struct.unpack('I', buffer[offset:offset+4])
    return the_int

#FAIL
#def decompress(firstbytes, saveStream):
#    """Decompress the presumed save file stream. First 4 bytes are already consumed, so we take those in as firstbytes parameter"""
#    #inputStream = StringIO.StringIO()
#    #outputStream = StringIO.StringIO()
#    #inputStream, outputStream = os.pipe()
#    #process = subprocess.Popen(['./blast'], stdin=inputStream,stdout=outputStream, shell=True)
#    process = subprocess.Popen(['./blast'], stdin=subprocess.PIPE,stdout=saveStream, shell=True)
#    process.communicate(firstbytes)
#    process.communicate(saveStream)
#    #response = process.communicate(firstbytes)
#    #response += process.communicate(saveStream)
#    #return outputStream

def parse_save(saveFile):
    buffer = saveFile.read(4)
    if buffer <> 'CIV3':
        print "wah wah wah wahhhhhhhh."
        print "Stub. Provided stream not decompressed C3C save"
        return -1
    #print 'Using Horspool search to go to first WRLD section'
    wrldOffset = horspool.boyermoore_horspool(saveFile, "WRLD")
    #print wrldOffset
    game = Wrld(saveFile)
    saveFile.close()
    return game

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

def main():
    saveFile = open("gamesaves/unc-test.sav", 'rb')
    #saveFile = open("gamesaves/unc-lk151-650ad.sav", 'rb')
    game = parse_save(saveFile)

    print 'Printing something(s) from the class to ensure I have what I intended'
    #print game.name, game.length
    #print game.tile.pop()[1].length
    print "Map width:",game.width,"Map height:", game.height
    print game.Tiles.tile[0].info['terrain']
    print game.Tiles.tile[1000].info['terrain']
    #print game.tile[0].is_visible_to
    #max = len(game.tile)
    max = 10
    for i in range(max):
        print game.Tiles.tile[i].continent

if __name__=="__main__":
    main()
