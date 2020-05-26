-- To help me make documentation
-- function tree_values(t)
    --
-- end

function tree_values(t, n)
    local w = io.write
    -- default hierarchy n to 0 if not set
    n = n or 0
    -- indent string, length based on hierarchy
    local s = string.rep("  ", n) .. "- "
    -- loop through all key/value pairs of the table
    for k,v in pairs(t) do
        w(s..k)
        -- if value is a table, recurse and increase indentation
        if type(v) == "table" then
            if v[1] then
                if type(v[1]) == "table" then
                    w("[]\n")
                    tree_values(v[1],n+1)
                else
                    w("[] (", type(v[1]), ")\n")
                    --
                end
            else
                w("\n")
                tree_values(v,n+1)
            end
        elseif type(v) == "function" then
            w("()\n")
        else
            -- otherwise print the key/value pair which are usually "tagType" and "name"
            w("\n")
        end
    end
end

bic.load_default()
sav.load(install_path .. "/Saves/Auto/Conquests Autosave 4000 BC.SAV")

tables = {
    bic = bic,
    sav = sav,
    civ3 = civ3,
    game = game,
    wrld = wrld,
    tile = tile,
    bldg = bldg,
    city = city,
    diff = diff,
    eras = eras,
    govt = govt,
    lead = lead,
    prto = prto,
    race = race,
    tech = tech,
    unit = unit,
    wsiz = wsiz,
}

-- I don't understand why tables is an array of tables, but it is
for k, v in pairs{tables} do
    -- print("- " .. k)
    tree_values(v)
end
