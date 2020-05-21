-- This file is what I'm using in development to test the code I'm working on at the moment

function lpad(s, l, c)
    local res = string.rep(c or ' ', l - #s) .. s
    return res
end

function mass_scan()
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
        io.write(lpad(tostring(suede.city_count), 4))
        io.write(lpad(tostring(suede.unit_count), 5))
        io.write(' ', save_name,'\n')
    end
end

function do_other_stuff()
    sav.load(install_path .. "/Saves/Auto/Conquests Autosave 4000 BC.SAV")
end

mass_scan()
-- do_other_stuff()