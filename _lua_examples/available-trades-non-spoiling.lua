-- Shows available tech trades in-game, in csv format

relative_save = "/Saves/Auto/Conquests Autosave 4000 BC.SAV"

player = 1

bic.load_default()
sav.load(install_path .. relative_save)

print("\"Civ\", ", "\"Tech to Buy\", ", "\"Tech to Sell\"")

-- Assumes Lua 1-based index and not civ3-native 0-based
function has_prereqs(civ, tech_id)
    if tech[tech_id].eras_id > lead[civ].eras_id then
        return false
    end
    for _, v in ipairs(tech[tech_id].prereq_tech_ids) do
        if bit32.band(game.tech_civ_bitmask[v+1], bit32.lshift(1, civ - 1)) == 0 then
            return false
        end
    end
    return true
end

for k, v in ipairs(lead) do
    if v.race_id > 0 and k - 1 ~= player then
        if v.contact_with[player + 1] ~= 0 then
            io.write("\"", race[v.race_id + 1].name, "\"")
            if v.at_war[player + 1] == 0 or v.will_talk_to[player + 1] == 0 then
                local tech_to_buy = {}
                local tech_to_sell = {}
                local player_mask = bit32.lshift(1, player)
                local civ_mask = bit32.lshift(1, k - 1)
                local civ = k
                for k, v in ipairs(game.tech_civ_bitmask) do
                    if (bit32.band(v, player_mask) == 0) ~= (bit32.band(v, civ_mask) == 0) then
                        if bit32.band(v, player_mask) == 0 then
                            if has_prereqs(player + 1, k) then
                                table.insert(tech_to_buy, tech[k].name)
                            end
                        else
                            if has_prereqs(civ, k) then
                                table.insert(tech_to_sell, tech[k].name)
                            end
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