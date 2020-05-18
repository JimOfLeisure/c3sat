## luaciv3

An attempt to expose Civ III SAV and BIQ files in a Lua scripting environment.
I'll be using gopher-lua.

### Dev notes

- should probably return a new lua environment; maybe optionally inject into existing
- should be able to read files direcly (and eventually possibly write)
- `sav` and `bic` table variables for parseciv3-similar queries by section
- Perhaps `game`, `wrld`, `tile`, etc. table variables for processed data?
- no spoiler protection at first, but maybe add later, perhaps as function on table vars or global setting
- My LevelDB storage idea doesn't belong in this package; should go in the executable
- maybe start with `sav.load()` and registry path finder
