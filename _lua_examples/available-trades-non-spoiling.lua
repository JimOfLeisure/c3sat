-- Shows available tech trades in-game
-- Needs better formatting ... csv, plain text, html?

relative_save = "/Saves/Auto/Conquests Autosave 4000 BC.SAV"
relative_save = "/Saves/Cleopatra of the Egyptians, 2310 BC.SAV"

player = 1

bic.load_default()
sav.load(install_path .. relative_save)

print("Civ", "Tech to Buy", "Tech to Sell")

for k, v in ipairs(lead) do
    if v.race_id > 0 and k - 1 ~= player then
        if v.contact_with[player + 1] ~= 0 then
            io.write(race[v.race_id + 1].name)
            if v.at_war[player + 1] == 0 or v.will_talk_to[player + 1] == 0 then
                -- print(' trade')
                local tech_to_buy = {}
                local tech_to_sell = {}
                local player_mask = bit32.lshift(1, player)
                local civ_mask = bit32.lshift(1, k - 1)
                for k, v in ipairs(game.tech_civ_bitmask) do
                    if (bit32.band(v, player_mask) == 0) ~= (bit32.band(v, civ_mask) == 0) then
                        if bit32.band(v, player_mask) == 0 then
                            --
                            table.insert(tech_to_buy, tech[k].name)
                        else
                            --
                            table.insert(tech_to_sell, tech[k].name)
                        end
                    end
                end
                io.write(", \"")
                for k, v in ipairs(tech_to_buy) do
                    io.write(v)
                    if k ~= #tech_to_buy then
                        io.write(', ')
                    end
                end
                io.write("\", \"")
                for k, v in ipairs(tech_to_sell) do
                    io.write(v)
                    if k ~= #tech_to_sell then
                        io.write(', ')
                    end
                end
                io.write("\"")
            else
                io.write("\"envoy refused\", \"envoy refused\"")
            end
            io.write("\n")
        end
    end
end