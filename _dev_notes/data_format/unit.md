# UNIT

2020-05-21 - Trying for the first time to decode UNITs.

- GAME offset 28 is an integer with the global unit count
- but this doesn't always exactly match the number of UNIT sections - when is a UNIT not a unit? (so far I've only seen more UNITs than unit count)
- LEAD also has per-civ unit counts at offset 368, an int
- Unit seems to consistenly have a 0x1d8 (472) size with 536 bytes between UNIT headers (includes IDLS section per UNIT)
- There do seem to be "extra" bytes after IDLS and before the next UNIT - Antal1987 dumps say it's 12 ints, and that looks about right
