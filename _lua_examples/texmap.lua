relative_save = "/Saves/Auto/Conquests Autosave 4000 BC.SAV"

-- 1 fits better on screens, 2 looks closer to right aspect ratio
-- 3 seems to behave oddly, but the aspect ratio is closer to right
text_tile_width = 1

water_tile = string.rep("~", text_tile_width*2)
land_tile = string.rep("█", text_tile_width*2)
indent = string.rep(" ", text_tile_width)

sav.load(install_path .. relative_save)

for k, v in ipairs(tile) do
    -- newline for end of map row
    if (k - 1) % (tile.width / 2) == 0 then
        io.write("\n")
        -- indent odd map rows
        if math.floor((k - 1) / (tile.width / 2)) % 2 == 1 then
            io.write("  ")
        end
    end
    if v.base_terrain > 10 then
        io.write(water_tile)
    else
        io.write(land_tile)
    end
end
io.write("\n")
