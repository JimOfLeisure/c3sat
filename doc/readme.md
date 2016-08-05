August 2016: During the Go rewrite, taking more specific notes. They might as well go here.

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
- "BICQVER#" C3C or "BICXVER#" PTW