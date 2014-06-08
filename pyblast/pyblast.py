#!/usr/bin/env python
# -*- coding: latin-1 -*-

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

# NOTICE: I am using blast.c from zlib's contrib/blast by Mark Adler and
# Ben Rudiak-Gold's report on his reverse-engineering of DCL as references
# in writing this module.
# ref: https://github.com/madler/zlib/tree/master/contrib/blast
# ref: https://groups.google.com/forum/#!msg/comp.compression/M5P064or93o/W1ca1-ad6kgJ
# None of this Python code was written by Mark Adler, but I will be attempting
# to heavily borrow from his C code. If it is possible to copy, paste and modify I
# will do so. That said, I will include his copyright and license notice here:

# /* blast.h -- interface for blast.c
#   Copyright (C) 2003, 2012 Mark Adler
#   version 1.2, 24 Oct 2012
# 
#   This software is provided 'as-is', without any express or implied
#   warranty.  In no event will the author be held liable for any damages
#   arising from the use of this software.
# 
#   Permission is granted to anyone to use this software for any purpose,
#   including commercial applications, and to alter it and redistribute it
#   freely, subject to the following restrictions:
# 
#   1. The origin of this software must not be misrepresented; you must not
#      claim that you wrote the original software. If you use this software
#      in a product, an acknowledgment in the product documentation would be
#      appreciated but is not required.
#   2. Altered source versions must be plainly marked as such, and must not be
#      misrepresented as being the original software.
#   3. This notice may not be removed or altered from any source distribution.
# 
#   Mark Adler    madler@alumni.caltech.edu
#  */

# NOTICE: I will restate that this is at most a heavy modifiction and port of
# Mark Adler's C code into Python by Jim Nelson. At least it is a from-scratch
# attempt at decompressing save game files using Mark Adler's code and
# Ben Rudiak-Gold's report as references.

# My aim is not for a general-purpose decoder. My aim is to decompress a
# stream from a file or url


def main():
    print "hi"

if __name__=="__main__":
    main()

