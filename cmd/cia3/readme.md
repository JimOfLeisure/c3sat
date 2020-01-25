Note: Must run `pkger -o /cmd/cia3` before building cmd/cia3 to create pkged.go
in this directory.

Also, pass the `-ldflags="-H windowsgui"` flag to `go build` to avoid having a
console window open while the GUI is open with Windows.
