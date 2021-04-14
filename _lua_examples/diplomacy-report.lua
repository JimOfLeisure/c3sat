-- A summary of diplomatic information in CSV format

-- Unlikely to have any contacts in 4000 BC, so change this
relative_save = "/Saves/Auto/Conquests Autosave 1300 BC.SAV"

player = 1

bic.load_default()
sav.load(install_path .. relative_save)

-- Takes civ as a Lua 1-based index to lead
function is_alive(civ)
    if lead[civ].city_count == 0 and lead[civ].unit_count == 0 then
        return false
    end
    return true
end

io.write("\"Civ\", \"Relationship\", \"Will Talk\", \"Government\", \"Era\", \"Cities\", \"Gold\"\n")

for k, v in ipairs(lead) do
    if v.race_id > 0 and k - 1 ~= player and is_alive(k) then
        if v.contact_with[player + 1] ~= 0 then
            io.write("\"", race[v.race_id + 1].name, "\", ")
            if v.at_war[player + 1] == 0 then
                io.write("\"Peace\", \"Yes\"")
            else
                if  v.will_talk_to[player + 1] == 0 then
                    io.write("\"WAR\", \"Yes\"")
                else
                    io.write("\"WAR\", \"NO\"")
                end
            end
            io.write(", \"", govt[v.govt_id + 1].name, "\", ")
            io.write(", \"", eras[v.eras_id + 1].name, "\", ")
            io.write(tostring(v.city_count), "\", ")
            io.write(tostring(v.gold), "\", ")
            io.write("\n")
        end
    end
end
