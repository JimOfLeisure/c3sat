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
    foo = get_savs({install_path .. "/Saves/Auto", install_path .. "/Saves"})
    for _, v in pairs(foo) do
        -- print(v)
        sav.load(v)
        -- print(sav.dump())
        -- print(civ3.always26)
        -- print(civ3.maybe_version_minor)
        -- print(civ3.maybe_version_major)
        -- for k, v in pairs(civ3) do
        --     print(k, v)
        -- end
        -- print(tile.width)
        -- print(tile.height)
        if suede.unit_sections ~= suede.unit_count then
            io.write(lpad(tostring(suede.city_count), 4))
            io.write(lpad(tostring(suede.unit_count), 5))
            io.write(lpad(tostring(suede.unit_sections), 6))
            io.write(' ', save_name,'\n')
        end
    end
end

function do_other_stuff()
    bic.load_default()
    -- sav.load(install_path .. "/Saves/Auto/Conquests Autosave 4000 BC.SAV")
    -- sav.load(install_path .. "/Saves/nice start Lincoln of the Americans, 4000 BC.SAV")
    sav.load(install_path .. "/Saves/Cleopatra of the Egyptians, 2310 BC.SAV")
    print(prto.dump)
    for k, v in pairs(prto) do
        print(k,v)
    end
    for k, v in ipairs(prto) do
        print("---", k)
        for kk, vv in pairs(v) do
            print(kk,vv)
        end
    end
end

-- mass_scan()
do_other_stuff()