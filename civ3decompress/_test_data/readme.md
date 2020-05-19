# Test data

AI.bin is just 8 bytes (hex 00 04 82 24 25 c7 80 7f) which is a compressed
"AIAIAIAIAIAIA" as described in
<https://groups.google.com/forum/#!msg/comp.compression/M5P064or93o/W1ca1-ad6kgJ>
. I used it to test when developing the decompressor, and decompress_test.go
now uses it.
