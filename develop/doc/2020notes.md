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

## LEAD
  - int32 length (offsets from below start after length)
  - 0x00: int32 player order? player index? (as expected)
  - 0x04: int32 race ID, -1 if not playing
  - 0x08: int32 starts game at 0
  - 0x0c: int32 starts game at 0 - Power (in F8 histograph) ✓
  - 0x10: int32 starts game at -1, then appears to be count from player# to 0 for non-barbs, in reverse order of index (additional: -1 until first city founded?)
  - 0x14: int32; 2 for AI, 3 for human player?
  - 0x18: int32 0 in early game
  - 0x1c: int32 0 in early game
  - 0x20: int32 -1 in early game (Golden age end?)
  - 0x24: int32 4 in early game (status?)
    - Noticed it changed to 12 for me. saw it a turn or two after making peace. Wonder if this is bit mask and 0x8 is declared war on someone?
    - ~~the next few bytes make me think there's a byte or char here somewhere; need to look w/hex dump~~ Still looks int32-aligned, perhaps these are encoded gold or some other encoded value (gold count is protected from easy hex editing)
  - 0x28: int32 encoded gold? but only lsb seems to change - maybe byte array? "lsb" seems to increment slowly
  - 0x2c: int32 encoded gold? but only lsb seems to change - maybe byte array? "lsb" seems to increment slowly
  - 0x30: possibly near start of byte or int16 array
  - 0x84 : int32 - mobilization level (?)
  - 0x88 : int32 - government type (?)
  - 0x8c : int32 - # of map tiles discovered (?)
  - 0x91-ish : This seems to occasionaly get civ name strings, but I think it's a bug and data should be ints of some length
  - 0xdc : int32 - ~~culture?~~ no, but have only seen it increment so far, but believe it must go down as I saw a 0
  - 0xe4 : int32 - ~~culture?~~ no, but have only seen it increment so far, but believe it must go down as I saw a 0
  - 0xec : ~~int32 - military unit count? or garrison count?~~ No, I think this part of an array, maybe int16s
  - 0xf2 (not int32 boundary, part of int32[60]) : started incrementing 1 per turn around turn 10/11?
    - I seem to be 2 ahead of 4 other civs, and two are still at 0
    - Then it didn't change for anyone
  - 0xfc : int32 - era?
  - 0x12c : 00 to 01 2950bc but don't know why
  - 0x132 : 00 to 01 2950bc but don't know why
  - 0x188 : int32 - tax luxury slider (0..10)
  - 0x18c : int32 - tax science slider (0..10)
  - 0x190 : int32 - tax cash (inferred) slider (0..10)

### some sort of diplomatic stuff, perhaps?

  - 0x1e0 : player 8 went from 0 to 1 when I declared war on player 8
    - same for player 7 when I dow'ed player 7
  - 0x200 : player 5 went from 0 to 2 when they declared war on me
    - player 2 went from 0 to 1 after making a demand and backing down when I refused
  - 0x320 : 00 to 01 when player 6 dow'ed after demand refusal
  - 0x328 : 00 to 01 when player 6 dow'ed after demand refusal
  - 0x338 : player 7 00 to 01 when I dow'ed
  - 0x41c : went from 00 to 01, not sure why
    - player 5 went 0 to 1 when dow'ing me
  - 0xb18 : player 8 decremented 0x17 to 0x16 during war w/me, unsure if related
    - decremented to 0x15 to 0x14, now will talk, unsure if related
    - not related to war or peace w/me
    - noticed player 7 decrementing after dow...maybe turns left on research? wild guess
    - player 5 went 0 to 20 after dow'ing me
    - player 2 went 1 to 20 after demanding money and backing down when I refused
  - 0xb98 : int32? - (array, 0 if player willing to talk?) player 8 went from 5 ~~(furious?)~~ to 4 ~~(annoyed?)~~ during war but won't speak
    - went to 3 next turn but still no speak, and annoyed
    - went to 2 next turn but still no speak, and annoyed
    - when went to 0, will speak!
    - player 8 00 to 07 when made peace after war, they're annoyed. Peace treaty related?
    - player 7 00 to 07 when I declared war. not peace/war-related?
    - player 5 00 to 08 when they declared war
  - 0xc20 : int32 - player 7 went from 03 to 06 when I declared war
  - 0xc28 : int32 - player 7 went from 03 to 01 when I declared war
  - 0xc98 : int32 - player 8 went from 0 to 0xffffffe2 (-30) when I declared war on player 8, they refuse to speak
### war
  - 0xd14 : Guessing this is always 01 for war vs barbs, presuming this is start of byte array for war.
  - 0xd15 : player 8 went from 0 to 1 when I declared war on player 8, they refuse to speak
    - player 8 went 1 to 0 when made peace
    - player 7 went from 0 to 1 when I declared war
    - player 5 went from 0 to 1 when they declared war
  - 0xd18 : went from 0 to 1 when I declared war on player 7
  - 0xd19 : went from 0 to 1 when player 5 declared war on me
  - 0xd1c : went from 0 to 1 when I declared war on player 8
    - and back to 0 when made peace with player 8
  - 0xd98 : player 8 went from 1 to 0 when I declared war on player 8, they refuse to speak
    - player 8 went from 1 to 0 when I declared war, refuse to speak
### contact

Bit flags?
```
0x1 = Contact
0x2 = ? unit in sight? Never spoke? Got this after saving after seeing new contact between turns
    seems to mean have contact but never actually spoke (spoke & saved game immediately, changed to just 0x1 from 3)
0x4 = ? Got this, 0x8, and 0x10 along with 0x1 after making peace and the civ having two archers left in my territory - possibly indicates can demand withdraw or declare; or might mean offensive unit
0x8 = their foreign unit in civ territory, seen when my scout is in their territory. same if I have a warrior in their territory
0x10 = ? Got this, 0x8, and 0x10 along with 0x1 after making peace and the civ having two archers left in my territory - possibly indicates can demand withdraw or declare; or might mean offensive unit
```
  - 0xe94 : presumed start of int32 contact array. This is 0 for barb player
  - 0xe98: player 5 went from 00 to 03 when I made contact with player 5 ("cautious" towards me? doesn't seem to line up with a bit flag for player 1)
    - player 8 also went 00 to 03 when met, and they are cautious
    - player 3 also went 00 to 03 when met, and they are cautious ("cautious" is not in the BIQ)
    - player 2 00 to 01 when met, and they are annoyed
    - player 7 00 to 03 when met, and they are polite
  - 0xe9c: went from 00 to 01 when I made contact with player 2
  - 0xea0: went from 00 to 01 when I made contact with player 3
  - 0xea8: went from 00 to 01 when I made contact with player 5
  - 0xeb0: went from 00 to 01 when I made contact with player 7
  - 0xeb4: went from 00 to 01 when I made contact with player 8

### ?

  - 0x11a0 & 0x11a4 ints radically changed on turn 10 when borders expanded; also happened to the AI civs same location same time, but no making sense of the numbers yet
    - also changed on turn 11 for all civs
    - and 12, and continually, apparently
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
  79..100 int field_13022];
  91 ~~101~~ int Armies_Count; ✓?
  92 ~~102~~ int Unit_Count; ✓
  93 ~~103~~ int Military_Units_Count; ✓
  94 ~~104~~ int Cities_Count; ✓
  105 int field_198;
  106 int field_19C;
  107 int field_1A0;
  98 ~~108~~ int Tax_Luxury; ✓
  100 ~~109~~ int Tax_Cash; ✓
  99 ~~110~~ int Tax_Science; ✓
  111..846 int field_1B0[736];
  0xd14 : char At_War[32];
  char field_D50[32];
  char field_D70[32];
  int field_D90[72];
  0xe94 : int Contacts[32];
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