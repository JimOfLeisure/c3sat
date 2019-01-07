## Civ3 Show-And-Tell

Civ3 Show-And-Tell reads `sav` and `biq` files from Civilization III Conquests and Civilization III Complete v1.22. It can:

- Show world seed and map settings needed to regenerate the same map
- Decompress a `sav` or `biq` file
- Hex dump the file data, automatically decompressing if needed
- Return data based on GraphQL queries

Other features are in development such as:

- Generating maps viewable in a web browser
- Reports on other game information

### Use

`civ3sat.exe --path <file> <command>`

This example runs the `seed` command against a save file. To generate the same
map you must use all the choices in the "choose your world" screen as shown in
the Choice column when starting a new game. The Result column shows the end
result which is only different if the original generator chose random.

Note that world size behaves differently. If a map was generated with random
world size, and you want to recreate it, you have to choose the resulting map
size or else the map may be different. 

    >civ3sat.exe seed --path "C:\Program Files (x86)\Steam\steamapps\common\Sid Meier's Civilization III Complete\Conquests\Saves\Auto\Conquests Autosave 4000 BC.SAV"

    Setting         Choice          Result
    World Seed      156059380
    World Size      Huge
    Barbarians      Random          No Barbarians
    Land Mass       Random          Archipelago
    Water Coverage  Random          70% Water
    Climate         Random          Arid
    Temperature     Random          Warm
    Age             Random          3 Billion

### Commands

- `seed` - Returns world map info as shown above.
- `decompress` - Writes a decompressed version of the file as `out.sav` in the current working directory.
- `hexdump` - Prints a hex dump of all game data to the console. If the file is compressed, it will automatically decompress first. 
- `graphql <query>` - Executes a GraphQL query
- `api` - Starts http server with GraphQL API at http://127.0.0.1/graphql

### GraphQL

When running the `api` command, andy GraphQL client can be used against http://127.0.0.1/graphql , or you can browse to that URL in the browser and it will load the Playground GraphQL client in the browser. You can also execute command-line quereies with the `graphql` command but will need to escape double quotes.

Queries defined:

- Direct data queries with section header and ordinal, offset from start of section and count of values
  - `bytes` - Returns byte array, assumes all bytes are unsinged
  - `int16s` - Returns int16 array, assumes all int16s are signed
  - `int32s` - Returns int16 array, assumes all int32s are signed
  - `hexString` - Like bytes but returns hex string, e.g. "0100FFFF"
  - `base64` - Like bytes but base64-encoded
- `civ3` - The first interpreted query returning named values. Use the GraphQL client or example below to see the available values

Example queries:

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
