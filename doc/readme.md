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
- int, always 1. Count of BICs? Count of whatever the next blob is?
- length, always 720 / 0x2d0
- length bytes of data
    - 0x00 - 0xff - four ints?
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
- "BLDG" in scenario-based saves, "GAME" in epic saves
    - There must be an indicator in the file previous to this to indicate what data is present/missing