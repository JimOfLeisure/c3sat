2021-11-05: I recently got the idea to show heat/value maps of map bytes. It will show 2d spatial shapes much better. Going through the TILEs looking for interesting patterns.

Offsets in this doc are from the first byte after the first TILE as 0x00.

0x04 - Rivers ~~Definite clumping pattern on land, most likely terrain related.~~
0x05 - Owned tile / culture borders (player number)

0x06 - all 179s? Middle gray / middle value
0x07 - all 9s

  *Possibly an int32 here*
0x08 - Resource graphics tile id, including volcanoes. Spotted differing numbers for volcanoes appearing differently
0x09 - Also associated with resources, but many zero out compared to 0x08. (Mayve all zeroes? Need to check.)
10 - Also resource-associated, but all zeroes?
11 - Also resource-associated, but all zeroes?

12 - Unit ID, I think. Almost certainly unit-correlated. Need to check to see if they are all unique
13
14
15 - \^ part of a unit ID

16 - Base terrain graphics tile index - This looks really funky in the value-as-texture view
17 - Base terrain graphics file ID

18 through 23 - All zeroes? (Possibly pre-Conquests data?)

24 - all 255 - From other notes, believed to be barb info, perhaps tribe ID
25 - all 255

26 - City ID
27 - All zeroes? Maybe part of city ID, but I think there may be 256 limit on cities

28 - all 255 - Believed to be Colony ID from other notes?
29 - all 255

30 - Continent ID
31 - All zeroes?

32 - All 6s

33 - Secret code? This seems to be zeroes except for horizontal data near the top of map. It's what I think it would look like if I stuffed data into a save with my idea for C3X. It doesn't really look like textual or structured data, or any reasonable collection of numbers, though. Not even particularly aligned in any way.

34 - 255s
35 - 255s

36 - Razed city! All 0s, but found a 1, and it's ruins in-game. Not seeing other values, so can rule out pollution. More 1s found for razed cities.

37 - 0s
38 - 0s
39 - 0s - possible msb of razed city int32? Doesn't really make sense

40 through 43 - "TILE"
44 through 47 = 12 (section length)

48 - Improvements flag (includes pollution which is 0x40 flag)

49 - 0s
50 - 0s
51 - 0s
52 - 0s

53 - Terrain (2 nybbles, ID referencdes to BIC)

58 - appears to be bit flag for start locations, bonus grasslands, forests chopped, etc. Bit flat 0x02 may be tile workable / available? See "Civ3C-SaveFileNotes" for more
59 - more bit flags; in fact I think the chopped one is 0x10 here.

60 through 63 - "TILE"
64 through 67 - 4 (32-bit length)

68 through 71 - 0s

72 through 75 - "TILE"
76 through 29 - 128 (32-bit length)