# CITY

2020-05-23 - Trying for the first time to decode CITYs.

- There appear to be usually 91 "CITY" section headers per city!!!
  - A couple of outliers with custom bic files downloaded from elsewhere
  - At least one was originally a PTW save, and the global city count was wrong, so the offset was wrong for PTW
  - The other was a Rise of Rome Conquest; the global city count seems plausible, and the "CITY"/ city count ratio seems steady at 33
  - Upon further investigation, the number of building types appears to represent that number of "CITY" headers, and there are 83 in default game
- Apparent structure from looking at one save (default epic game, autosave 2270 BC w/19 cities)
    - CITY 0x88 (length/count)
        - CITY 0x10
        - CITY 0x24 - all 0xff in my example save
        - CITY 0xa4
        - CITY 0x94 Name string at 0x4
        - POPD 0x08 - potential CTZN count at 0x8
        - 3* CTZN 0x12c (assume one per citizen)
        - BINF 0x04 - I've always presumed this was "bin food", but that seems to be elsewhere. This might be "building info" as it has a value of 83
        - Lots of non-sectioned data, seems to alternate between 0 and -1 - Think
         this is "building info", seems to be 83 triplets of ints, or maybe
         shorts or something, but lenth 12 bytes * count. *** Looks like it's
         3 ints: year, civ ID, culture
        - BITM 0x28
        - DATE 0x54
        - CITY 0x08
        - CITY 0x04 - List header with count (83 - 0x53)
          - 83* CITY 0x04 There are 83 building types in the default game; wonder if this is a list of maybe per-building data like date built?
        - CTPG 0x4
        - CTPG 0x10
        - CITY 0x04

