
-- left padding for easy terminal viewing
function lpad(s, l, c)
    local res = string.rep(c or ' ', l - #s) .. s
    return res
end

    bic.load_default()
    io.write("\"cities\",\"ocn\",\"barb camps\",\"barb wars\",\"barb horses\",\"cities/ocn\"\n")
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
        local barb_camp = 0
        for _, v in ipairs(tile) do
            if bit32.band(v.improvements, 0x80) ~= 0 then
                barb_camp = barb_camp + 1
            end
        end
        local ocn = wsiz[wrld.wsiz_id + 1].ocn
        ocn = math.floor(ocn * diff[game.diff_id].pct_optimal_cities / 100)
        io.write(lpad(tostring(game.city_count), 4))
        io.write(',', lpad(tostring(ocn), 4))
        io.write(',', lpad(tostring(barb_camp), 4))
        io.write(',', lpad(tostring(barb_warrior), 4))
        io.write(',', lpad(tostring(barb_horseman), 4))
        io.write(', ', string.format("%.2f", game.city_count / ocn))
        -- io.write(', \"', save_name,'\"')
        io.write('\n')
    end
