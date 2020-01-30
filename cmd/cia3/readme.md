Note: Must run `pkger -o /cmd/cia3` before building cmd/cia3 to create pkged.go
in this directory.

Also, pass the `-ldflags="-H windowsgui"` flag to `go build` to avoid having a
console window open while the GUI is open with Windows.

### Bundled Javascript packages

I've included some js bundles in this repo. I'm currently using them for analyzing diffs between hex dumps in civs.html. I don't have a one-step build script; I just copied them over manually.

- [jsdiff](https://github.com/kpdecker/jsdiff) - html/diff has the dist/ from this package
- [diff2html](https://github.com/rtfpessoa/diff2html) - html/diff2html has some of the bundles/ files from this package