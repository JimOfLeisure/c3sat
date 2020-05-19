bic.load_default()
sav.load(install_path .. "/Saves/Auto/Conquests Autosave 4000 BC.SAV")
print(tile.width)
print(tile.height)
for k, v in ipairs(tile) do
    -- print(v.terrain)
    -- print(v.base_terrain)
    print(v.overlay_terrain)
end


foo = get_savs({install_path .. "/Saves/Auto"})
for _, v in pairs(foo) do
    print(v)
    sav.load(v)
    -- print(sav.dump())
    -- print(civ3.always26)
    -- print(civ3.maybe_version_minor)
    -- print(civ3.maybe_version_major)
    -- for k, v in pairs(civ3) do
    --     print(k, v)
    -- end
    print(tile.width)
    print(tile.height)
end