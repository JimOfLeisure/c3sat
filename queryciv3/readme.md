# queryciv3

QueryCiv3 implements a GraphQL API into a Civ3 SAV & BIQ file combination.

It is currently used by both /cmd/c3sat and /cmd/cia3. It originally handled
the http server and still has that functionality, but I've tried to make it more
an http handler provider.

The following information is copied from the old root readme.md and it not
fully up-to-date, but it's not wrong, either.

More recent examples and use can be found in /cmd/cia3/html/cia3.js , although
that's a rather ugly long file in early 2020.

## Queries defined

- Direct data queries with section header and ordinal, offset from start of section and count of values
  - `bytes` - Returns byte array, assumes all bytes are unsinged
  - `int16s` - Returns int16 array, assumes all int16s are signed
  - `int32s` - Returns int16 array, assumes all int32s are signed
  - `hexString` - Like bytes but returns hex string, e.g. "0100FFFF"
  - `base64` - Like bytes but base64-encoded
- `civ3` - The first interpreted query returning named values. Use the GraphQL client or example below to see the available values

## Example queries

- Starting locations of players

        {
            int32s(section: "WRLD", nth: 2, offset: 36, count: 32)
        }

- Tile's trade network ID by civ; nth should be a multiple of 4

        {
            int16s(section: "TILE", nth: 4, offset: 26, count: 32)
        }

- Get map generation values

        {
            civ3 {
                worldSeed
                size
                barbariansFinal
                landMassFinal
                oceanCoverageFinal
                climateFinal
                temperatureFinal
                ageFinal
            }
        }
