## Civ3 Show-And-Tell

Civ3 Show-And-Tell reads `sav` and `biq` files from Civilization III Conquests and Civilization III Complete v1.22. It can:

- Show world seed and map settings needed to regenerate the same map
- Decompress a `sav` or `biq` file
- Hex dump the file data, automatically decompressing if needed
- Return data based on GraphQL queries

Other features are in development such as:

- Generating maps viewable in a web browser
- Reports on other game information

### Binaries

- [Windows 386](http://lib.bigmoneyjim.com/civ3sat/0.3.1/windows-386/civ3sat.exe)
    - MD5 D97D5C5F69EB7D6609EDEAC60FE7880B
    - SHA1 4DB6DDB1B98A904ADD0CCCFFFCA9C505F5357758
    - SHA256 45AAE4A14A900CE4A3C9E93957215575449DB281D8AC05680A0DF7D65DC5F961
- [Linux 386](http://lib.bigmoneyjim.com/civ3sat/0.3.1/linux-386/civ3sat)
    - MD5 EDD80CA17BB2C72F319835D5BAFA4669
    - SHA1 914D7FB868C14118BF33C488C01CF9FEC692D418
    - SHA256 147587A323441F424DEF73A8DE7E7C54F8DAEA5F07ED6282D983DD1C451DCB41
- [Darwin/Mac 386](http://lib.bigmoneyjim.com/civ3sat/0.3.1/darwin-386/civ3sat)
    - MD5 8B5A5742EC2FB2E8A33D01FC988DB30C
    - SHA1 00091CAB05EEA445CBFE329EAEC5C428CEAD366B
    - SHA256 8D6DCB256082B62B3339C722E238AFA2594E6001E213C710999E21F056A40D70

### Use

`civ3sat.exe --path <file> <command>`

The following example runs the `seed` command against a save file. To generate the same
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
