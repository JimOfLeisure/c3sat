One day it would be nice to re-organize these notes, but for now (January 2020) I'll just add some new (re-?)discoveries.
I'm focusing only on v1.22 (final) Conquests/Complete. The info may or may not be similar for earlier patches or versions.

I'm now pretty sure that an entire BIQ is embedded in the save file when using any custom scenario/conquest.
Otherwise it refers to the conquests.biq file. I'm not sure if sections of the biq can be eliminated and fall
back to the default biq, or if the whole biq has to be present. I'm also not sure if custom maps in BIQ files
remain in the BIQ section. They certainly wouldn't be needed once the SAV is created because it has its own
copy of the map.

I believe the second GAME section is the beginning of the per-game save data, and the rest of the file is part of that
top-level structure. Note that if "GAME" is part of any of the BIQ strings this could throw off seek-by-text
algorithms.

```
'CIV3'
    'BIC '
        int32?
        int32?
        path name support files location?
        path name biq location?
            'BICX', 'BICQ', or similar; the beginning of the entire BIQ file
            ...
            'GAME'
    'GAME'
        ...
        everything else
```

WRLD, TILE, and CONT Seem pretty straightforward and covered in other notes.

After CONT seems to be an unnamed section that is an int32 array of resource counts. The length is the number of
resource types are defined (GOOD from the BIQ). If this length is repeated elsewhere I haven't found it yet,
but I would look in that GAME section starting the SAV for such counts.

Then LEAD begins. It appears every save has 32 LEADs. Each begins with a length in bytes, but there is more
data after that. I believe they are int32 arrays whose lenths are based possibly on **types** of resources, units,
cty improvements (buildings), techs, and maybe some other stuff. Then a couple of ESPN sections and CULT, each
with a length in bytes. And then there is another unnamed/count-not-included int32 array.

[Antal1987's dumps](https://github.com/myjimnelson/C3CPatchFramework/blob/master/Civ3/Leader.h) may be instructive in helping to look what data is there.

Backing up to the BIQ's RACE section: RACE appears to be what I call a basic list section. It starts with an
int32 count (always 32?) of 'races'/civs. Each list item has the correct byte count. However, the first
data structure in the item is a list of cities, and the number of cities is inconsistent from civ to civ.
It begins with a count of cities, each of which seems to be 32-character 0-terminated strings (Windows-1252 encoding). The downside of this is that the other
civ data is not the same offset from each item starting offset, so you'll have to parse the city list
to find the offset where the other civ data begins.

The game data refers to the BIQ data by index, and the BIQ is where all the strings are.
I think there might be some text files used for human language translations, but I'm not sure.
