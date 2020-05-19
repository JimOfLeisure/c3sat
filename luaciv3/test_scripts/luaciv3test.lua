bic.load_default()
sav.load(install_path .. "/Saves/Auto/Conquests Autosave 4000 BC.SAV")
print(sav.dump())
print(bic.dump())
foo = get_savs({install_path .. "/Saves/Auto"})
for _, v in pairs(foo) do
    print(v)
    sav.load(v)
    print(sav.dump())
end