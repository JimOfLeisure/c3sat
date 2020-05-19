# Civ3 Show-And-Tell

This repo contains Go code to decompress and read save and BIQ files from Civilization III Conquests and Civilization III Complete v1.22.

There are three executables, c3sat, cia3, and clua.

The libraries implement a GraphQL API local web server to query the data, and the cia3 executable includes embedded JavaScript to
query the saves and present useful information.

The luaciv3 libary implements a Lua environment with modules to access Civ3 save file data. See /luaciv3/readme.md for more detail.

## Civ Intelligence Agency III (CIA3)

Civ Intelligence Agency III is intended to be a non-spoiling, non-cheating game assistant for single-human-player games of Sid Meier's Civilization III Complete and/or Conquests. Multiplayer non-spoiling support may be possible in the future.

CIA3 watches save game folder locations for new .SAV files, then reads the file and updates the information. Every time you save your game or begin a new turn, a new save or autosave will be generated and trigger CIA3 to update its info. It does not update mid-turn unless you save your game during the turn.

CIA3 currently reads the Windows registry to determine the default save and autosave folders for Civ3 and watches those. In the future I'll add the ability to watch specified folders and perhaps open files on demand.

The information provided by CIA3 is either available in-game, was available earlier in-game or during map creation, or accepted as not cheating by Civ III competitions. For example, the player may not be able to determine in-game if forests have been previously chopped for shield production on newly-revealed map tiles, but previous assistants (CivAssist II) have shown this information and have been allowed in competitive games. Another example is that previous tools allowed exact map tile counts when at least 80% of the map is revealed to allow for determining how far the player is from winning by domination.

Previous assistants such as <a href="https://forums.civfanatics.com/resources/crpsuite-2-11-0.16266/">CRpSuite's MapStat</a> <a href="https://forums.civfanatics.com/resources/crpsuite-2-11-0.16266/">(download page)</a> and <a href="https://forums.civfanatics.com/threads/civassist-ii.118540/">CivAsist II</a> <a href="https://forums.civfanatics.com/resources/civassist-ii.21/">(download page)</a> worked well for many years, but neither of these seem to work in Windows 10.

I, Jim Nelson, or <a href="https://forums.civfanatics.com/members/puppeteer.36357/">Puppeteer on CivFanatics Forums</a>, have been on-and-off since 2013 been working on a save game reader for other purposes with <a href="https://forums.civfanatics.com/threads/civ3-show-and-tell.493582/">Civ3 Show-And-Tell</a> (C3SAT). Until 2020 I explicitly did <strong>not</strong> want to recreate a game assistant, but now that the others aren't working and seem to be abandoned by their creators, I started working on CIA3 from the C3SAT code base. CIA3 is a different product from C3SAT, but they share a common code base for most functions.

<a href="https://forums.civfanatics.com/threads/cia3-civ-intelligence-agency-iii.656876/">Release thread on CivFanatics Forums</a> where new releases are announced. <a href="https://forums.civfanatics.com/threads/civ3-show-and-tell.493582/">Discussion thread on CivFanatics Forums</a> where I babble about development progress. CIA3 discussion begins on <a href="https://forums.civfanatics.com/threads/civ3-show-and-tell.493582/page-10#post-15638589">page 10</a>

### Current Features

- Available tech trades
- Map highlighting forest chopped squares
- Contact, will-talk and war/peace status for opponent civs
- Display of map generation settings, world seed, and difficulty
- Should work with custom scenarios/conquests

### Known Issues

- Eliminated civs still show up in lists
- At-war won't-talk civs will disappear from tech trades table until they're willing to talk

## Civ3 Show-And-Tell (C3SAT)

v0.4.1 Note: C3SAT has not changed functionality since v0.4.0, except that the new GraphQL queries would be available in API and query modes, but I don't
intend to compile and release binaries, so either use CIA3's binary, grab the binary release from v0.4.0, or build from source.

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

## Civ3 Lua (clua)

`clua` behaves similarly to the standard `lua` executable, but it's Lua version
5.1, and it has modules to access Civ3 save file data. See /luaciv3/readme.md,
/\_lua\_examples, and Lua 5.1 guides for scripting.

It is the newest and least-complete of the programs here, and very little is
anywhere near complete.