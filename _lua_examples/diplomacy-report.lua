-- A summary of diplomatic information in CSV format

-- Unlikely to have any contacts in 4000 BC, so change this
relative_save = "/Saves/Auto/Conquests Autosave 4000 BC.SAV"
relative_save = "/Saves/Cleopatra of the Egyptians, 2310 BC.SAV"

player = 1

bic.load_default()
sav.load(install_path .. relative_save)

io.write("\"Civ\", \"Relationship\", \"Will Talk\", \"Government\", \"Era\", \"Cities\"\n")

for k, v in ipairs(lead) do
    if v.race_id > 0 and k - 1 ~= player then
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
            io.write(tostring(v.city_count))
            io.write("\n")
        end
    end
end
