relative_save = "/Saves/Auto/Conquests Autosave 4000 BC.SAV"

sav.load(install_path .. relative_save)

for k, v in ipairs(tile) do
    -- newline for end of map row
    if (k - 1) % (tile.width / 2) == 0 then
        io.write("\n")
        -- indent odd map rows
        if math.floor((k - 1) / (tile.width / 2)) % 2 == 1 then
            io.write(" ")
        end
    end
    -- tilde for water, █ for land
    if v.base_terrain > 10 then
        io.write("~~~~")
    else
        io.write("████")
    end
end
io.write("\n")
