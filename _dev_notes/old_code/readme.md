# Old code

Stuff I'm keeping around for reference and possible updating later.

## Files and folders

- d3js/ folder - I believe this expects the data in an older API format and
doesn't work with the current format. I don't recall when it worked; it may be
from before I had the executable with the API function. I think it's from the
old python version, and this read a json output file of that code. It should be
trivial but tedious to update it to work the the current REST API.
- ApolytonToGoStruct.ps1 - A few times I've tried to rigidly define the data
structures so they can just read the file as a serialized data dump which I
believe it is. Most of the time I think that's too much trouble because it's so
hard to get right without the original source code. But this seems to be a
PowerShell script to create Go data structure definitions based on the
file format research published on Apolyton, and I presume the much-modified
result is what is currently in the /parseciv3/types.go file in this repo.
- Build.ps1 - This looks like a PowerShell build and test script I made and
used at some point. I really should learn cmake, and nowadays I'd make some
_test.go files instead of an external script.
- mapdefs.svg - When I was making SVG-based maps, I think this was the file
with all the shapes defined, and I'd include this from the map SVG. I don't
think this has been used since the project was in Python.

## Other git points of interest in this project

This project started out in Python. It had a very brief early detour into C#
(<https://github.com/myjimnelson/c3satcs>), I think went back to Python, and
was redone in Go several years ago.

### Python code

The referenced to the latest Python code (I think) is in the
[abandonedpydevelop](https://github.com/myjimnelson/c3sat/tree/abandonedpydevelop)
tag.

The Python code was mostly focused on interpreting and generating maps from
save files. it produced output files which were at various times svg, json,
or html. Well, and text dumps in the early days when I was learning the very
basics of the file format.

There was also a brief working attempt at a uWGSI web app that would take save
file uploads and process them into an HTML/SVG map for online viewing.

### Webview branch

The [webview branch](https://github.com/myjimnelson/c3sat/tree/webview) tried
to use github.com/zserge/webview as a UI instead of
v0.4.x's fyne-io with link for any browser. In Jan 2020 webview still used
IE11 on Windows but had a webview-x branch that used Edge on windows. Webview-x
using Edge was supposed to be ready and in the master branch by March or April. If I
keep using a web UI I may revisit this branch and see about launching straight
into a 'native' browser window.