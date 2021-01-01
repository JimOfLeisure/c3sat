This is an attempt to export the Civ3 decompressor as a C shared libary

## Build commands

### Windows

Requires MinGW gcc or TDM-GCC

`go build -o decompressciv3.dll -buildmode=c-shared`

### MacOS

`go build -o decompressciv3.dylib -buildmode=c-shared`

## Linux

`go build -o decompressciv3.so -buildmode=c-shared`
