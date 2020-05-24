-- This file is what I'm using in development to test the code I'm working on at the moment

function lpad(s, l, c)
    local res = string.rep(c or ' ', l - #s) .. s
    return res
end

-- pass it a table and optionally hierarchy number
function printTables(t, n)
    -- default hierarchy n to 0 if not set
    n = n or 0
    -- indent string, length based on hierarchy
    local s = string.rep("  ", n)
    -- loop through all key/value pairs of the table
    for k,v in pairs(t) do
        -- if value is a table, recurse and increase indentation
        if type(v) == "table" then
            print(s..k..":")
            printNbtTables(v,n+1)
        else
            -- otherwise print the key/value pair which are usually "tagType" and "name"
            print(s..k, v)
        end
    end
end

function mass_scan()
    bic.load_default()
    foo = get_savs({
        install_path .. "/Saves/Auto",
        install_path .. "/Saves",
        install_path .. "/Saves\\ancient-saves\\CivIII-Conquests-v1.22-Saves",
        install_path .. "/Saves\\ancient-saves\\CivIII-Conquests-v1.22-Saves\\Auto",
        install_path .. "/Saves/2017-2018",
        -- this has just one file, a PTW save
        -- install_path .. "/Saves/civfan",
        install_path .. "/Saves/huge france autosaves",
        install_path .. "/Saves/Twitch",
        install_path .. "/Saves/YouTube",
    })
    for _, v in pairs(foo) do
        print(v)
        sav.load(v)
        io.write(lpad(tostring(game.city_count), 4))
        io.write(', \"', save_name,'\"\n')
    end
end

function do_other_stuff()
    bic.load_default()
    -- sav.load(install_path .. "/Saves/Auto/Conquests Autosave 4000 BC.SAV")
    -- sav.load(install_path .. "/Saves/nice start Lincoln of the Americans, 4000 BC.SAV")
    -- sav.load(install_path .. "/Saves/Cleopatra of the Egyptians, 2310 BC.SAV")
    sav.load(install_path .. "/Saves/Auto/Conquests Autosave 450 BC.SAV")
    -- print(prto.dump)
    for k, v in ipairs(city) do
        print("-----", k,v)
        for kk, vv in pairs(v) do
            -- print(kk,vv)
            io.write(lpad(tostring(kk), 25))
            io.write(" ", lpad(tostring(vv), 6))
            io.write("\n")
        end
        -- for k, v in ipairs(v.binf) do
        --     for kk, vv in pairs(v) do
        --         print(kk,vv)            
        --     end
        -- end
end
    -- for k, v in pairs(tech) do
    --     print("---", k)
    --     for kk, vv in pairs(v) do
    --         print(kk,vv)            
    --     end
    --     for kk, vv in pairs(v.prereq_tech_ids) do
    --         print(kk,vv)
    --     end
    -- end
    -- for k, v in ipairs(game.tech_civ_bitmask) do
    --     print(tech[k].name,v)
    -- end
    -- print(bit32.band(15,2))
end

function textmap()
    bic.load_default()
    relative_save = "/Saves/Cleopatra of the Egyptians, 2310 BC.SAV"

    -- 1 fits better on screens, 2 looks closer to right aspect ratio
    -- 3 seems to behave oddly, but the aspect ratio is closer to right
    text_tile_width = 2
    
    water_tile = string.rep("~", text_tile_width*2)
    land_tile = string.rep("█", text_tile_width*2)
    indent = string.rep(" ", text_tile_width)
    
    sav.load(install_path .. relative_save)
    
    for k, v in ipairs(tile) do
        -- newline for end of map row
        if (k - 1) % (tile.width / 2) == 0 then
            io.write("\n")
            -- indent odd map rows
            if math.floor((k - 1) / (tile.width / 2)) % 2 == 1 then
                io.write(indent)
            end
        end
        foo = v.improvements
        if bit32.band(foo, 0x80) ~= 0 then
        -- if foo ~= 0 then
                -- io.write(lpad(tostring(v.continent_id), text_tile_width*2))
            io.write(lpad("o", text_tile_width*2))
        else
            -- io.write(lpad(" ", text_tile_width*2))
            if v.base_terrain > 10 then
                io.write(water_tile)
            else
                io.write(land_tile)
            end
            end
    end
    io.write("\n")
end

mass_scan()
-- do_other_stuff()
-- textmap()