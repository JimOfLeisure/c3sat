
-- left padding for easy terminal viewing
function lpad(s, l, c)
    local res = string.rep(c or ' ', l - #s) .. s
    return res
end

    bic.load_default()
    io.write("\"cities\",\"bar wars\",\"barb horses\",\"save name\"\n")
    civ3_save_files = get_savs({
        install_path .. "/Saves/Auto",
        install_path .. "/Saves",
    })
    for _, v in pairs(civ3_save_files) do
        sav.load(v)
        local barb_horseman = 0
        local barb_warrior = 0
        for _, v in ipairs(unit) do
            if v.civ_id == 0 and v.prto_id == 6 then
                barb_warrior = barb_warrior + 1
            end
            if v.civ_id == 0 and v.prto_id == 11 then
                barb_horseman = barb_horseman + 1
            end
        end
        io.write(lpad(tostring(game.city_count), 4))
        io.write(',', lpad(tostring(barb_warrior), 4))
        io.write(',', lpad(tostring(barb_horseman), 4))
        io.write(', \"', save_name,'\"\n')
    end
