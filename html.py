#!/usr/bin/env python

import tileonly
import datetime


def main():
    outputhtmlpath = '/usr/share/nginx/www/map.html'
    htmlhead = 'html/head'
    htmltail = 'html/tail'

    game = tileonly.parse_save()

    write = open(outputhtmlpath, 'w')
    head = open(htmlhead, 'r')
    tail = open(htmltail, 'r')

    write.write(head.read())
    write.write(str(datetime.datetime.now()))
    #write.write(game.html_out())
    #write.write(game.html_fake_iso())
    write.write(game.table_out())
    write.write(tail.read())

main()
