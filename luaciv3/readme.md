# luaciv3

LuaCiv3 exposes Civ III SAV and BIQ file data in a Lua scripting environment.

Currently everything is read-only. Write functionality *may* be added in the
future, but I'm not sure it's worth the trouble.

## Global variables & functions

- `install_path` -  If found in the Windows registry, this is the path to the
Civ III Conquests (or Complete) install directory. Of course, it can be overridden in Lua.
- `get_savs({dir})` or `get_savs({dir1, dir2, dir3})` - Given an array of
directory paths, returns an array of all '.sav' files in those directories.
(Not currently recursive.)
- `sav_path` - The full path name of the loaded SAV file
- `sav_name` - Just the file name of the loaded SAV file
- `bic_path` - The full path name of the default BIC file
- `bic_name` - Just the file name of the default BIC file

## Raw data modules

- `sav` - Information and functions about the raw SAV file
  - `sav.load(path)` - Given the path name, loads a save file into memory and
  populates the other modules' data.
  - `sav.dump()` - Currently hex dumps the first 30 bytes of the save file. May
  be expanded or eliminated in the future.
- `bic` - Information and functions about the raw BIC/BIX/BIQ
  - `bic.load_default(path)` - Loads a BIC in memory to be used as the default
  BIC when a custom BIC is not present in a SAV file. If no path provided, it
  will try to load `install_path .. "/conquests.biq"`.
  - `bic.dump()` - Currently hex dumps the first 256 bytes of the default BIC in
  memory. May be expanded or eliminated in the future.

## Processed data modules

- `civ3` - Exposes some data in the file header or "CIV3 section". Not really
useful as I'm not sure the data structure is right. But it was a decent
proof-of-concept first module.
- `tile` - Exposes map data from the "TILE sections" and some helper data from
elsewhere. See /\_lua\_examples/textmap.lua to see how to iterate the tile map
and organize the output tiles.
  - `tile.width` - Map width in tiles. Note that there are only `tile.width / 2`
  tiles on each row due to the way civ3 does its coordinates.
  - `tile.height` - Map height in tiles.
  - `tile[1]`, `tile[2]`, etc.. - Iterated with `for k, v in ipairs(tile) do; <...> end`. Note that in-game tile offsets are 0-indexed and lua is 1-indexed, so you'll need to add or subtract one from the index when converting from civ3 offsets to lua offsets or vice versa.
    - `tile[n].terrain` - The terrain byte for the tile. This isn't useful unprocessed and may be removed in the future.
    - `tile[n].base_terrain` - The low nybble of terrain. Indicates the base terrain type.
    - `tile[n].overlay_terrain` - The high nybble of terrain. Indicates the overlay terrain (hill, mountain, forest, volcano, flood plain, marsh, or jungle) else is just a duplicate of the base terrain.

## Automated dump of data types

- bic
  - load_default()
  - dump()
- sav
  - load()
  - dump()
- civ3
  - always26
  - maybe_version_minor
  - maybe_version_major
  - gobbledegook_1
  - gobbledegook_2
  - gobbledegook_3
  - gobbledegook_4
- game
  - diff_id
  - unit_count
  - city_count
  - tech_civ_bitmask[] (number)
- wrld
  - wsiz_id
- tile[]
  - terrain
  - base_terrain
  - overlay_terrain
  - continent_id
  - improvements
  - overlay
- bldg[]
  - dump
- city[]
  - name
  - id
  - x
  - y
  - lead_id
  - improvements_maintenance
  - stored_food
  - stored_shields
  - ctzn_count
  - ctzn[]
    - dump
  - binf[]
    - date
    - lead_id
    - culture
  - date
  - city[] (number)
- diff[]
  - name
  - content_citizens
  - max_anarchy_turns
  - defense_land_units
  - offense_land_units
  - start_units_1
  - start_units_2
  - add_free_support
  - bonus_each_city
  - barb_attack_bonus
  - cost_factor
  - pct_optimal_cities
  - ai_trade_rate
  - corruption
  - quelled_citizens
- eras[]
  - name
- govt[]
  - name
- lead[]
  - city_count
  - unit_count
  - player_number
  - race_id
  - govt_id
  - mobilization_level
  - tiles_discovered
  - eras_id
  - research_beakers
  - current_research_id
  - current_research_turns
  - future_techs_count
  - armies_count
  - military_unit_count
  - at_war[] (number)
  - will_talk_to[] (number)
  - contact_with[] (number)
- prto[]
  - name
  - attack
  - defense
  - move
  - cost
  - transport
- race[]
  - city_names[] (string)
  - great_leader_names[] (string)
  - leader_name
  - leader_title
  - adjective
  - name
  - object_noun
- tech[]
  - name
  - eras_id
  - prereq_tech_ids
- unit[]
  - id
  - x
  - y
  - prev_x
  - prev_y
  - civ_id
  - race_id
  - prto_id
- wsiz[]
  - name
  - ocn
  - tech_rate
  - width
  - dist_between_civs
  - num_civs
  - height
