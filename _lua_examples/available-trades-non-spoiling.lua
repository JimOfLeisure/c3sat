-- File is under development.
-- The goal is to show non-spoiling tech trade opportunities

relative_save = "/Saves/Auto/Conquests Autosave 4000 BC.SAV"
relative_save = "/Saves/Cleopatra of the Egyptians, 2310 BC.SAV"

player = 1

bic.load_default()
sav.load(install_path .. relative_save)

print("Civ", "Tech to Buy", "Tech to Sell")

for k, v in pairs(lead) do
    if v.race_id > 0 and k - 1 ~= player then
        print(race[v.race_id + 1].name)
    end
end