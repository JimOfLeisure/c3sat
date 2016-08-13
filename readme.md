## Civ3 Show-And-Tell

Civ3 Show-And-Tell reads `sav` and `biq` files from Civilization III Conquests and Civilization III Complete v1.22. It can:

- Show world seed and map settings needed to regenerate the same map
- Decompress a `sav` or `biq` file
- Hex dump the file data, automatically decompressing if needed

Other features are in development such as:

- Generating maps viewable in a web browser
- Reports on other game information

### Use

`civ3sat.exe <command> <file>`

This example runs the `seed` command against a save file. To generate the same
map you must use all the choices in the "choose your world" screen as shown in
the Choice column when starting a new game. The Result column shows the end
result which is only different if the original generator chose random.

Note that world size behaves differently. If a map was generated with random
world size, and you want to recreate it, you have to choose the resulting map
size or else the map may be different. 

    >civ3sat.exe seed "C:\Program Files (x86)\Steam\steamapps\common\Sid Meier's Civilization III Complete\Conquests\Saves\Auto\Conquests Autosave 4000 BC.SAV"

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

- `seed <file>` - Returns world map info as shown above.
- `decompress <file>` - Writes a decompressed version of the file as `out.sav` in the current working directory.
- `hexdump <file>` - Prints a hex dump of all game data to the console. If the file is compressed, it will automatically decompress first. 
