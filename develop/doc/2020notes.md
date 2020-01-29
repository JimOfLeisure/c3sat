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

- LEAD
  - int32 length (offsets from below start after length)
  - 0x00: int32 player order? player index? (as expected)
  - 0x04: int32 race ID, -1 if not playing
  - 0x08: int32 starts game at 0
  - 0x0c: int32 starts game at 0 - Power? (in F8 histograph)
  - 0x10: int32 starts game at -1, then appears to be count from player# to 0 for non-barbs, in reverse order of index (additional: -1 until first city founded?)
  - 0x14: int32; 2 for AI, 3 for human player?
  - 0x18: int32 0 in early game
  - 0x1c: int32 0 in early game
  - 0x20: int32 -1 in early game (Golden age end?)
  - 0x24: int32 4 in early game (status?)
    - ~~the next few bytes make me think there's a byte or char here somewhere; need to look w/hex dump~~ Still looks int32-aligned, perhaps these are encoded gold or some other encoded value (gold count is protected from easy hex editing)
  - 0x28: int32 encoded gold? but only lsb seems to change - maybe byte array? "lsb" seems to increment slowly
  - 0x2c: int32 encoded gold? but only lsb seems to change - maybe byte array? "lsb" seems to increment slowly
  - 0x30: possibly near start of byte or int16 array
  - 0x84 : int32 - mobilization level (?)
  - 0x88 : int32 - government type (?)
  - 0x8c : int32 - # of map tiles discovered (?)
  - 0x91-ish : This seems to occasionaly get civ name strings, but I think it's a bug and data should be ints of some length
  - 0xdc : int32 - culture?
  - 0xe4 : int32 - culture?
  - 0xec : int32 - military unit count? or garrison count?
  - 0xfc : int32 - era?
  - 0xea8: went from 00 to 01 when I made contact with player 5
  - 0xe98: went from 00 to 03 when I made contact with player 5 ("cautious" towards me? doesn't seem to line up with a bit flag for player 1)
  - ... end of LEAD length
- int32 array(s)
- ESPN
- ESPN
- CULT
- int32 array(s)

[Antal1987's dumps](https://github.com/myjimnelson/C3CPatchFramework/blob/master/Civ3/Leader.h) may be instructive in helping to look what data is there.

Backing up to the BIQ's RACE section: RACE appears to be what I call a basic list section.

- It starts with an int32 count (always 32?) of 'races'/civs. Each list item has the correct byte count.
- However, the first data structure in the item is a list of cities, and the number of cities is inconsistent from civ to civ.
- It begins with a count of cities, each of which seems to be 24-character 0-terminated strings (Windows-1252 encoding). The downside of this is that the other civ data is not the same offset from each item starting offset, so you'll have to parse the city list (and more) to find the offset where the other civ data begins.
- Following the city list is an int32 count of 16-character military great leader names list.
- Then leader name string 32 chars
- Then leader title string 24 chars
- e.g. "RACE_Romans" string 32 chars
- e.g. "Roman" adjective string 40 chars?
- Civ name string 40 chars
- e.g. "Romans" object noun string 40 chars
- start of several strings referencing 8 flc files, strings 256+ chars each
- a few int32s I think
- int32 count and list of scientific great leader names 
- end of RACE section item

The game data refers to the BIQ data by index, and the BIQ is where all the strings are.
I think there might be some text files used for human language translations, but I'm not sure.

## from Antal1987's Lead.h

With notes added.

```
  int field_4[6]; - length?
  0 int ID; ✓ 
  1 int RaceID; ✓
  2, 3 int field_24[2];
  4 int CapitalID;
  5 int field_30;
  6 int field_34;
  7 int field_38;
  8 int Golden_Age_End;
  9 int Status;
  10 int Gold_Decrement;
  11 int Gold_Encoded;
  12..32 int field_4C[21];
  33 int GovenmentType;
  34 int Mobilization_Level;
  35 int Tiles_Discovered;
  36..49 int field_AC[14];
  50 int field_E4;
  51..53 int field_E8[3];
  54 int Era; ✓
  55 int Research_Bulbs; ✓
  56 int Current_Research_ID; ✓
  57 int Current_Research_Turns; ✓
  58 int Future_Techs_Count; ✓?
  59..78 __int16 AI_Strategy_Unit_Counts[20];
  79..100 int field_130[22];
  101 int Armies_Count;
  102 int Unit_Count;
  103 int Military_Units_Count;
  104 int Cities_Count;
  105 int field_198;
  106 int field_19C;
  107 int field_1A0;
  108 int Tax_Luxury;
  109 int Tax_Cash;
  110 int Tax_Science;
  111..846 int field_1B0[736];
  char At_War[32];
  char field_D50[32];
  char field_D70[32];
  int field_D90[72];
  int Contacts[32];
  int Relation_Treaties[32];
  int Military_Allies[32];
  int Trade_Embargos[32];
  int field_10B0[18];
  int Color_Table_ID;
  int field_10FC;
  int field_1100[7];
  int field_111C[36];
  int field_11AC[8];
  int field_11CC;
  int field_11D0[252];
  int field_15C0;
  int field_15C4;
  int field_15C8;
  int field_15CC;
  int Science_Age_Status;
  int Science_Age_End;
  int field_15D8;
  __int16 *Improvement_Counts;
  int field_15E0;
  int Improvements_Counts;
  int *Small_Wonders;
  int field_15EC;
  __int16 *Units_Counts;
  int field_15F4;
  int field_15F8;
  __int16 *Spaceship_Parts_Count;
  int *ContinentStat4;
  int ContinentStat3;
  int *ContinentCities;
  int ContinentStat2;
  int *ContinentStat1;
  byte *Available_Resources;
  byte *Available_Resources_Counts;
  class_Civ_Treaties Treaties[32];
  class_Culture Culture;
  class_Espionage Espionage_1;
  class_Espionage Espionage_2;
  int field_18C0[260];
  class_Leader_Data_10 Data_10_Array2[32];
  class_Leader_Data_10 Data_10_Array3[32];
  class_Hash_Table Auto_Improvements;
```

gql query:

```
{
  civs {
    raceId: int32s(offset:4, count: 1)
    governmentType: int32s(offset:132, count: 1)
    mobilizationLevel: int32s(offset:136, count: 1)
    tilesDiscovered: int32s(offset:140, count: 1)
    era: int32s(offset:252, count: 1)
    UNSUREresearchBulbs: int32s(offset:256, count: 1)
    UNSUREcurrentResearchId: int32s(offset:260, count: 1)
    UNSUREcurrentResearchTurns: int32s(offset:264, count: 1)
    UNSUREfutureTechsCount: int32s(offset:268, count: 1)
  }
}
```