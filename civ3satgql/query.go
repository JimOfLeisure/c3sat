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
				// var ok bool
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
				// mdata.tilesData = make([][]byte, mdata.tileCount())
				// Read raw tile data
				mdata.tilesOffset, err = SectionOffset("TILE", 1)
				if err != nil {
					return nil, err
				}
				//  TODO: figure out how to handle wrapping, including oddball settings like Y wrap or no X wrap
				// 		Because I realized minX and maxX are inadequate in the case of a partial world across the wrap boundary
				// var minX, minY, maxX, maxY int
				// if playerSpoilerMask != 0 {
				// 	minX = mdata.mapWidth - 1
				// 	minY = mdata.mapHeight - 1
				// } else {
				// 	maxX = mdata.mapWidth - 1
				// 	maxY = mdata.mapHeight - 1
				// }
				//  *** TODO: and actually I don't need to read tile data into a buffer because saveGame.data exists in package context; use math during tile generation
				// for i := 0; i < mdata.tileCount(); i++ {
				// 	mdata.tilesData[i] = saveGame.data[section+i*196 : section+i*196+196]
				// 	// minX/etc logic goes here, use playerSpoilerMask and tile offset 68 I think / first value of 4th TILE / TILE128
				// }
				// mdata.tileSetX = minX
				// mdata.tileSetY = minY
				// mdata.tileSetWidth = maxX - minX + 1
				// mdata.tileSetHeight = maxY - minY + 1
				mdata.tileSetX = 0
				mdata.tileSetY = 0
				mdata.tileSetWidth = mdata.mapWidth
				mdata.tileSetHeight = mdata.mapHeight
				// mdata.tileSetX = mdata.mapWidth/2 - 5
				// mdata.tileSetY = 15
				// mdata.tileSetWidth = mdata.mapWidth / 2
				// mdata.tileSetHeight = mdata.mapHeight - 15
				return mdata, nil
			},
		},
		"bytes": &graphql.Field{
			Type:        graphql.NewList(graphql.Int),
			Description: "Byte array",
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
				return saveGame.data[savSection+offset : savSection+offset+count], nil
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
				return hex.EncodeToString(saveGame.data[savSection+offset : savSection+offset+count]), nil
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
	},
})
