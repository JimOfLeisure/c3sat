package civ3satgql

import (
	"encoding/base64"
	"encoding/hex"

	"github.com/graphql-go/graphql"
)

var queryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"civ3": &graphql.Field{
			Type:        civ3Type,
			Description: "Civ3 save data",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				wrldSection, err := SectionOffset("WRLD", 1)
				if err != nil {
					return nil, err
				}
				return worldData{worldOffset: wrldSection}, nil
			},
		},
		"fullPath": &graphql.Field{
			Type:        graphql.String,
			Description: "Save file path",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return saveGame.path, nil
			},
		},
		"fileName": &graphql.Field{
			Type:        graphql.String,
			Description: "Save file name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return saveGame.fileName(), nil
			},
		},
		"map": &graphql.Field{
			Type:        mapType,
			Description: "Current Game Map",
			Args: graphql.FieldConfigArgument{
				"playerSpoilerMask": &graphql.ArgumentConfig{
					Type:        graphql.Int,
					Description: "Bitmask of map tile per-player spoilers; default is 0x2 to show first human player. Set to 0 to return all map tiles.",
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var mdata mapData
				if spoilerMask, ok := p.Args["playerSpoilerMask"].(int); ok {
					mdata.playerSpoilerMask = int32(spoilerMask)
				} else {
					mdata.playerSpoilerMask = int32(0x2)
				}
				// Get 2nd WRLD offset
				section, err := SectionOffset("WRLD", 2)
				if err != nil {
					return nil, err
				}
				// Reading 6 int32s at offset 8; first is height, last is Width
				intList := make([]int, 6)
				for i := 0; i < 6; i++ {
					intList[i] = ReadInt32(section+8+4*i, Signed)
				}
				mdata.mapHeight = intList[0]
				mdata.mapWidth = intList[5]
				mdata.tilesOffset, err = SectionOffset("TILE", 1)
				if err != nil {
					return nil, err
				}
				var mapRowLength = mdata.mapWidth / 2
				var mapTileCount = mapRowLength * mdata.mapHeight
				mdata.mapTileOffsets = make([]int, mapTileCount)
				var minY, maxY int
				var mapXVisible = make([]bool, mapRowLength)
				minY = mdata.mapHeight - 1
				for i := 0; i < mapTileCount; i++ {
					tileOffset := mdata.tilesOffset - 4 + (mdata.tileSetY+i/mapRowLength)*mapRowLength*tileBytes + (mdata.tileSetX+i%mapRowLength)*tileBytes
					if mdata.spoilerFree(tileOffset) {
						mdata.mapTileOffsets[i] = tileOffset
						x := i % mapRowLength
						mapXVisible[x] = true
						y := i / mapRowLength
						if y < minY {
							minY = y
						}
						if y > maxY {
							maxY = y
						}
					} else {
						// tile elements will return null if offset <= 0
						mdata.mapTileOffsets[i] = -1
					}
				}
				// Need to ensure the first Y row is even because of how each odd row shifts a half tile right
				if minY%2 != 0 {
					minY -= 1
				}
				// See if it makes sense to wrap around X
				var longestBlank int
				var blankLength int
				// loop through width twice to ensure find the longest run of blanks, if any
				for i := 0; i < mapRowLength*2; i++ {
					if mapXVisible[i%mapRowLength] {
						if blankLength > 0 {
							if blankLength > longestBlank {
								longestBlank = blankLength
								mdata.tileSetX = (i % mapRowLength) * 2
							}
							blankLength = 0
						}
					} else {
						blankLength++
					}
				}
				mdata.tileSetWidth = (mapRowLength - longestBlank) * 2
				mdata.tileSetHeight = maxY - minY + 1
				mdata.tileSetY = minY
				var tileSetRowLength = mdata.tileSetWidth / 2
				var tileSetCount = tileSetRowLength * mdata.tileSetHeight
				var minX = mdata.tileSetX / 2
				mdata.tileSetOffsets = make([]int, tileSetCount)
				for i := 0; i < tileSetCount; i++ {
					mdata.tileSetOffsets[i] = mdata.mapTileOffsets[(i/tileSetRowLength+minY)*mdata.mapWidth/2+(i%tileSetRowLength+minX)%(mdata.mapWidth/2)]
				}
				return mdata, nil
			},
		},
		"bytes": &graphql.Field{
			Type:        graphql.NewList(graphql.Int),
			Description: "Byte array",
			Args: graphql.FieldConfigArgument{
				"target": &graphql.ArgumentConfig{
					Type:         graphql.String,
					Description:  "Target scope of the query. Can be game, bic, or file (default)",
					DefaultValue: "file",
				},
				"section": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Four-character section name. e.g. TILE",
				},
				"nth": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.Int),
					Description: "e.g. 2 for the second named section instance",
				},
				"offset": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.Int),
					Description: "Offset from start of section",
				},
				"count": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.Int),
					Description: "Number of bytes to return",
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var target *saveGameType
				targetArg, _ := p.Args["target"].(string)
				section, _ := p.Args["section"].(string)
				nth, _ := p.Args["nth"].(int)
				offset, _ := p.Args["offset"].(int)
				count, _ := p.Args["count"].(int)
				switch targetArg {
				case "game":
					target = &currentGame
				case "bic":
					target = &currentBic
				default:
					target = &saveGame
				}
				savSection, err := target.sectionOffset(section, nth)
				if err != nil {
					return nil, err
				}
				return target.data[savSection+offset : savSection+offset+count], nil
			},
		},
		"base64": &graphql.Field{
			Type:        graphql.String,
			Description: "Base64-encoded byte array",
			Args: graphql.FieldConfigArgument{
				"section": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Four-character section name. e.g. TILE",
				},
				"nth": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.Int),
					Description: "e.g. 2 for the second named section instance",
				},
				"offset": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.Int),
					Description: "Offset from start of section",
				},
				"count": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.Int),
					Description: "Number of bytes to return",
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				section, _ := p.Args["section"].(string)
				nth, _ := p.Args["nth"].(int)
				offset, _ := p.Args["offset"].(int)
				count, _ := p.Args["count"].(int)
				savSection, err := SectionOffset(section, nth)
				if err != nil {
					return nil, err
				}
				return base64.StdEncoding.EncodeToString(saveGame.data[savSection+offset : savSection+offset+count]), nil
			},
		},
		"hexString": &graphql.Field{
			Type:        graphql.String,
			Description: "Byte array in hex string format",
			Args: graphql.FieldConfigArgument{
				"section": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Four-character section name. e.g. TILE",
				},
				"nth": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.Int),
					Description: "e.g. 2 for the second named section instance",
				},
				"offset": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.Int),
					Description: "Offset from start of section",
				},
				"count": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.Int),
					Description: "Number of bytes to return",
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				section, _ := p.Args["section"].(string)
				nth, _ := p.Args["nth"].(int)
				offset, _ := p.Args["offset"].(int)
				count, _ := p.Args["count"].(int)
				savSection, err := SectionOffset(section, nth)
				if err != nil {
					return nil, err
				}
				return hex.EncodeToString(saveGame.data[savSection+offset : savSection+offset+count]), nil
			},
		},
		"hexDump": &graphql.Field{
			Type:        graphql.String,
			Description: "Hex dump of data",
			Args: graphql.FieldConfigArgument{
				"section": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Four-character section name. e.g. TILE",
				},
				"nth": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.Int),
					Description: "e.g. 2 for the second named section instance",
				},
				"offset": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.Int),
					Description: "Offset from start of section",
				},
				"count": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.Int),
					Description: "Number of bytes to return",
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				section, _ := p.Args["section"].(string)
				nth, _ := p.Args["nth"].(int)
				offset, _ := p.Args["offset"].(int)
				count, _ := p.Args["count"].(int)
				savSection, err := SectionOffset(section, nth)
				if err != nil {
					return nil, err
				}
				return hex.Dump(saveGame.data[savSection+offset : savSection+offset+count]), nil
			},
		},
		"int16s": &graphql.Field{
			Type:        graphql.NewList(graphql.Int),
			Description: "Int16 array",
			Args: graphql.FieldConfigArgument{
				"section": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Four-character section name. e.g. TILE",
				},
				"nth": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.Int),
					Description: "e.g. 2 for the second named section instance",
				},
				"offset": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.Int),
					Description: "Offset from start of section",
				},
				"count": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.Int),
					Description: "Number of int16s to return",
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				section, _ := p.Args["section"].(string)
				nth, _ := p.Args["nth"].(int)
				offset, _ := p.Args["offset"].(int)
				count, _ := p.Args["count"].(int)
				savSection, err := SectionOffset(section, nth)
				if err != nil {
					return nil, err
				}
				intList := make([]int, count)
				for i := 0; i < count; i++ {
					intList[i] = ReadInt16(savSection+offset+2*i, Signed)
				}
				return intList, nil
			},
		},
		"int32s": &graphql.Field{
			Type:        graphql.NewList(graphql.Int),
			Description: "Int32 array",
			Args: graphql.FieldConfigArgument{
				"section": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Four-character section name. e.g. TILE",
				},
				"nth": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.Int),
					Description: "e.g. 2 for the second named section instance",
				},
				"offset": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.Int),
					Description: "Offset from start of section",
				},
				"count": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.Int),
					Description: "Number of int32s to return",
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				section, _ := p.Args["section"].(string)
				nth, _ := p.Args["nth"].(int)
				offset, _ := p.Args["offset"].(int)
				count, _ := p.Args["count"].(int)
				savSection, err := SectionOffset(section, nth)
				if err != nil {
					return nil, err
				}
				intList := make([]int, count)
				for i := 0; i < count; i++ {
					intList[i] = ReadInt32(savSection+offset+4*i, Signed)
				}
				return intList, nil
			},
		},
		"allStrings": &graphql.Field{
			Type:        graphql.NewList(graphql.String),
			Description: "All ASCII strings four bytes or longer",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var i, count, offset int
				var output = make([]string, 0)
				for i < len(saveGame.data) {
					if saveGame.data[i] < 0x20 || saveGame.data[i] > 0x7F {
						if count > 3 {
							s := string(saveGame.data[offset:i])
							output = append(output, s)
						}
						count = 0
					} else {
						if count == 0 {
							offset = i
						}
						count++
					}
					i++
				}
				return output, nil
			},
		},
		"listSection": &graphql.Field{
			Type:        graphql.NewList(listSectionItem),
			Description: "A list section has a 4-byte count of list items, and each item has a 4-byte length",
			Args: graphql.FieldConfigArgument{
				"section": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Four-character section name. e.g. TILE",
				},
				"nth": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.Int),
					Description: "e.g. 2 for the second named section instance",
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				section, _ := p.Args["section"].(string)
				nth, _ := p.Args["nth"].(int)
				savSection, err := SectionOffset(section, nth)
				if err != nil {
					return nil, err
				}
				count := ReadInt32(savSection, Signed)
				output := make([]int, count)
				offset := 4
				for i := 0; i < count; i++ {
					output[i] = savSection + offset
					length := ReadInt32(savSection+offset, Signed)
					offset += 4 + length
				}
				return output, nil
			},
		},
		"civs": &graphql.Field{
			Type:        graphql.NewList(gameLeadSectionType),
			Description: "A list of 32 civilization/leader (LEAD) sections' data",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// Conquests saves appear to always have 32 LEAD sections
				const leadCount = 32
				output := make([]int, leadCount)
				for i := 0; i < leadCount; i++ {
					savSection, err := SectionOffset("LEAD", i+1)
					if err != nil {
						return nil, err
					}
					output[i] = savSection
				}
				return output, nil
			},
		},
	},
})
