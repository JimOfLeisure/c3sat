bic.load_default()
sav.load(install_path .. "/Saves/Auto/Conquests Autosave 4000 BC.SAV")
print(sav.dump())
print(bic.dump())
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
    for k, v in pairs(tile) do
        print(k, v)
    end
end