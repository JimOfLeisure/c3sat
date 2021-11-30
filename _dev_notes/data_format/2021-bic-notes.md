2021-11-30: I have Lua hooked up to my QueryCiv3 code and am looking for a "has custom rules" flag as some bics are just media reference, but they may also have custom rules and/or a custom map

civ3PTW folder has 87 bic/bix's

14 seem to have custom rules (BLDG and such)
69 seem to have maps
meaning 4 are media-only?

Output below. Offset 3 is the X or blank in the first 4 chars.
0x24 and 0x28 would appear to be LSBs for two int32s right before the first
string field.

The script looked for values that differ between bics in the first 32 bytes
and lists the values it finds for each offset that differs.

```
done    87
BLDG    14
WCHR    69
offset  3       
___
58      58
20      29



offset  24      
___
0b      58
02      26
04      1
03      2



offset  28      
___
12      49
05      6
0a      22
0d      3
09      2
01      1
08      2
06      1
07      1
```
