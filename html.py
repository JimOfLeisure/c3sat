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

import tileonly
import datetime


def main():
    """This module instantiates tileonly.parse_save() and writes an html file for the map"""
    outputhtmlpath = 'html/map.html'
    htmlhead = 'html/head'
    htmltail = 'html/tail'

    game = tileonly.parse_save()

    write = open(outputhtmlpath, 'w')
    head = open(htmlhead, 'r')
    tail = open(htmltail, 'r')

    write.write(head.read())
    write.write('<div>' + str(datetime.datetime.now()) + '</div>')
    #write.write(game.html_out())
    #write.write(game.html_fake_iso())
    #write.write(game.isometbrick_out())
    write.write('<div id="map">')
    write.write(game.svg_out())
    write.write(tail.read())
    write.write('</div>')

main()
