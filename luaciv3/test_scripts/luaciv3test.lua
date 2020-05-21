foo = get_savs({install_path .. "/Saves/Auto", install_path .. "/Saves"})
for _, v in pairs(foo) do
    -- print(v)
    sav.load(v)
    print(save_name)
    -- print(sav.dump())
    -- print(civ3.always26)
    -- print(civ3.maybe_version_minor)
    -- print(civ3.maybe_version_major)
    -- for k, v in pairs(civ3) do
    --     print(k, v)
    -- end
    -- print(tile.width)
    -- print(tile.height)
    print(suede.city_count)
end