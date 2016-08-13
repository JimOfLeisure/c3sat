August 2016: During the Go rewrite, taking more specific notes. They might as well go here.

## Start of SAV file

- "CIV3"
- 26 bytes of data that is the same for all my saves but seems to be different for other versions and GOTM, COTM and LK154/CCM, so perhaps created from BIC affects this
- "BIC "
- length
- length bytes of data. Seems to be 524 bytes in C3C and 525 bytes in PTW
    - 4 bytes in C3C, 5 bytes in PTW data
    - 256-byte (?) string relative path to BIC/X/Q resources
    - 4-byte zeroed int ?
    - 256-byte (?) string relative path to BIC/X/Q file
    - 4-byte zeroed int ?

## BIC section of SAV file. Same as BIC/X/Q file?

- "BICQVER#" C3C or "BICXVER#" PTW or "BIC VER#" vanilla
- int, always 1. Count of BICs? Count of whatever the next blob is? Version number as the previous bytes label it?
- length, always 720 / 0x2d0
- length bytes of data
    - 0x00 - 0xcdcdcdcd or 0
    - 0x04 - 0
    - 0x08 - BIC version
    - 0x0c - BIC version
    - 0x10 - string of well over 256 bytes - Description of scenario from BIQ
        - COTM0120_OPEN
        - COTM0121_OPEN
        - GOTM0149_OPEN
        - GOTM0150_OPEN
        - Mesopotamia is the "cradle of civilization" and was home to all seven of the "Great Wonders of the Ancient World." The map stretches from the mountains of Greece east to the hills of Persia and south to the Nile river. The game ends as soon as all seven great wonders have been built, or when a side amasses 5500 VP, or after 160 turns -- whichever happens first. Victory points awarded for completing wonders are doubled in this scenario
    - 0x290 - string - name of scenario
        - CCM
        - GOTM0150_OPEN
        - GOTM
        - COTM
        - Mesopotamia
    - All other seen bytes zeroes

### This part may not be present in epic game saves

There must be an indicator in the file previous to this to indicate what data is present/missing. I suspect it's in the first 16 bytes of the 720-byte chunk of data.

CivFanantic forum user Quintillus says that BIQ files with BLDG are due to custom rules, and WCHR is an optional custom map section.

Also, PTW and vanilla data structure sizes diverge beginning with the first BLDG

- "BLDG"
- int count
- count * 0x110 (272) byte records (length is different in PTW)
    - 0x44 - string - 32 bytes? display name of building
        - "Theory of Evolution"
    - 0x64 - string - 32 bytes? internal(?) name of building
        - "BLDG\_Theory\_of\_Evolution"
    - the rest seem to be ints
        - many -1
        - many 0
        - some other values
- "CTZN"
- int count
- count * 0x80 byte length laborer description
    - 0x00 - int value 0x7c
    - 0x04 - int, 1 for default laborer, 0 for the specialists
    - 0x08 - string display name
        - Entertainer
    - 0x28 - string internal name
        - CTZN_Entertainer
    - 0x48 - string plural
        - Entertainers
- "CULT"
- int count
- count * 0x5c byte length relative culture levels
    - 0x00 - int value 0x58
    - 0x04 - display text
        - impressed with
        - disdainful of
    - 0x44 - 6 * int
- "DIFF"
- int count
- count * 0x7c byte length difficulty levels
    - 0x00 - int value 0x78
    - 0x04 - display text
    - 0x44 - many ints
- "ERAS"
- int count
- count * 0x10c byte length eras
    - 0x00 - int value 0x108
    - 0x04 - display text
    - 0x44 - internal text
    - also names like "Wise-Men", "Alchemists", etc.
- "ESPN"
- int count
- count * 0xec byte length espionage actions
- "EXPR"
- int count
- count * 0x2c byte length unit experience levels

#### This may or may not be present

Think I saw flavor flags in the previous custom rules. Perhaps those can be checked. 

- "FLAV"
- int count (always 1?)
- count * flavor groups
    - int count flavors
    - count * 0x124 byte flavors

#### Back to present for all custom rules

- "GOOD"
- int count
- count * 0x5c byte goods/resources
- "GOVT"
- int count
- count * 0x23c byte governments
