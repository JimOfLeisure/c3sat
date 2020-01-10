package civ3satgql

import (
	"encoding/hex"

	"github.com/graphql-go/graphql"
)

type worldData struct {
	worldOffset int
}

type mapData struct {
	mapWidth          int
	mapHeight         int
	tileSetWidth      int
	tileSetHeight     int
	tileSetX          int
	tileSetY          int
	playerSpoilerMask int32
	tilesOffset       int
}

// no longer needed
// func (m *mapData) tileCount() int {
// 	return m.mapWidth * m.mapHeight
// }

var mapTileType = graphql.NewObject(graphql.ObjectConfig{
	Name: "tile",
	Fields: graphql.Fields{
		"foo": &graphql.Field{
			Type:        graphql.String,
			Description: "foo",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if offset, ok := p.Source.(int); ok {
					return string(saveGame.data[offset:offset+4]) + string(saveGame.data[offset+212:offset+216]), nil
				}
				return "foo", nil
			},
		},
		"hexTerrain": &graphql.Field{
			Type:        graphql.String,
			Description: "Byte value. High nybble is overlay, low nybble is base terrain",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if offset, ok := p.Source.(int); ok {
					return hex.EncodeToString(saveGame.data[offset+57 : offset+58]), nil
				}
				return "foo", nil
			},
		},
	},
})

var mapType = graphql.NewObject(graphql.ObjectConfig{
	Name: "map",
	Fields: graphql.Fields{
		"mapWidth": &graphql.Field{
			Type:        graphql.Int,
			Description: "Width of the game map in tiles",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if mdat, ok := p.Source.(mapData); ok {
					return mdat.mapWidth, nil
				}
				// TODO: better logic error handling?
				return -1, nil
			},
		},
		"mapHeight": &graphql.Field{
			Type:        graphql.Int,
			Description: "Height of the game map in tiles",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if mdat, ok := p.Source.(mapData); ok {
					return mdat.mapHeight, nil
				}
				// TODO: better logic error handling?
				return -1, nil
			},
		},
		"tileSetWidth": &graphql.Field{
			Type:        graphql.Int,
			Description: "Width of the currently visible map in tiles",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if mdat, ok := p.Source.(mapData); ok {
					return mdat.tileSetWidth, nil
				}
				// TODO: better logic error handling?
				return -1, nil
			},
		},
		"tileSetHeight": &graphql.Field{
			Type:        graphql.Int,
			Description: "Height of the currently visible map in tiles",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if mdat, ok := p.Source.(mapData); ok {
					return mdat.tileSetHeight, nil
				}
				// TODO: better logic error handling?
				return -1, nil
			},
		},
		"tileSetX": &graphql.Field{
			Type:        graphql.Int,
			Description: "World map X coordinate of top-left tile set tile",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if mdat, ok := p.Source.(mapData); ok {
					return mdat.tileSetX, nil
				}
				// TODO: better logic error handling?
				return -1, nil
			},
		},
		"tileSetY": &graphql.Field{
			Type:        graphql.Int,
			Description: "World map Y coordinate of top-left tile set tile",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if mdat, ok := p.Source.(mapData); ok {
					return mdat.tileSetY, nil
				}
				// TODO: better logic error handling?
				return -1, nil
			},
		},
		"tiles": &graphql.Field{
			Type:        graphql.NewList(mapTileType),
			Description: "Height of the currently visible map in tiles",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if mdat, ok := p.Source.(mapData); ok {
					const tileBytes = 212
					var tileRowLength = mdat.tileSetWidth / 2
					var mapRowLength = mdat.mapWidth / 2
					var tileCount = tileRowLength * mdat.tileSetHeight
					offsets := make([]int, tileCount)
					for i := 0; i < tileCount; i++ {
						offsets[i] = mdat.tilesOffset - 4 + (mdat.tileSetY+i/tileRowLength)*mapRowLength*tileBytes + (mdat.tileSetX+i%tileRowLength)*tileBytes
					}
					return offsets, nil
				}
				return nil, nil
			},
		},
	},
})

var civ3Type = graphql.NewObject(graphql.ObjectConfig{
	Name: "civ3",
	Fields: graphql.Fields{
		"worldSeed": &graphql.Field{
			Type:        graphql.Int,
			Description: "Random seed of random worlds",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if wdat, ok := p.Source.(worldData); ok {
					return ReadInt32(wdat.worldOffset+170, Signed), nil
				}
				return -1, nil
			},
		},
		"climate": &graphql.Field{
			Type: graphql.Int,
			// Description: "Random seed of random worlds",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if wdat, ok := p.Source.(worldData); ok {
					return ReadInt32(wdat.worldOffset+186, Signed), nil
				}
				return -1, nil
			},
		},
		"climateFinal": &graphql.Field{
			Type: graphql.Int,
			// Description: "Random seed of random worlds",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if wdat, ok := p.Source.(worldData); ok {
					return ReadInt32(wdat.worldOffset+190, Signed), nil
				}
				return -1, nil
			},
		},
		"barbarians": &graphql.Field{
			Type: graphql.Int,
			// Description: "Random seed of random worlds",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if wdat, ok := p.Source.(worldData); ok {
					return ReadInt32(wdat.worldOffset+194, Signed), nil
				}
				return -1, nil
			},
		},
		"barbariansFinal": &graphql.Field{
			Type: graphql.Int,
			// Description: "Random seed of random worlds",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if wdat, ok := p.Source.(worldData); ok {
					return ReadInt32(wdat.worldOffset+198, Signed), nil
				}
				return -1, nil
			},
		},
		"landMass": &graphql.Field{
			Type: graphql.Int,
			// Description: "Random seed of random worlds",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if wdat, ok := p.Source.(worldData); ok {
					return ReadInt32(wdat.worldOffset+202, Signed), nil
				}
				return -1, nil
			},
		},
		"landMassFinal": &graphql.Field{
			Type: graphql.Int,
			// Description: "Random seed of random worlds",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if wdat, ok := p.Source.(worldData); ok {
					return ReadInt32(wdat.worldOffset+206, Signed), nil
				}
				return -1, nil
			},
		},
		"oceanCoverage": &graphql.Field{
			Type: graphql.Int,
			// Description: "Random seed of random worlds",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if wdat, ok := p.Source.(worldData); ok {
					return ReadInt32(wdat.worldOffset+210, Signed), nil
				}
				return -1, nil
			},
		},
		"oceanCoverageFinal": &graphql.Field{
			Type: graphql.Int,
			// Description: "Random seed of random worlds",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if wdat, ok := p.Source.(worldData); ok {
					return ReadInt32(wdat.worldOffset+214, Signed), nil
				}
				return -1, nil
			},
		},
		"temperature": &graphql.Field{
			Type: graphql.Int,
			// Description: "Random seed of random worlds",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if wdat, ok := p.Source.(worldData); ok {
					return ReadInt32(wdat.worldOffset+218, Signed), nil
				}
				return -1, nil
			},
		},
		"temperatureFinal": &graphql.Field{
			Type: graphql.Int,
			// Description: "Random seed of random worlds",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if wdat, ok := p.Source.(worldData); ok {
					return ReadInt32(wdat.worldOffset+222, Signed), nil
				}
				return -1, nil
			},
		},
		"age": &graphql.Field{
			Type: graphql.Int,
			// Description: "Random seed of random worlds",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if wdat, ok := p.Source.(worldData); ok {
					return ReadInt32(wdat.worldOffset+226, Signed), nil
				}
				return -1, nil
			},
		},
		"ageFinal": &graphql.Field{
			Type: graphql.Int,
			// Description: "Random seed of random worlds",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if wdat, ok := p.Source.(worldData); ok {
					return ReadInt32(wdat.worldOffset+230, Signed), nil
				}
				return -1, nil
			},
		},
		"size": &graphql.Field{
			Type: graphql.Int,
			// Description: "Random seed of random worlds",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if wdat, ok := p.Source.(worldData); ok {
					return ReadInt32(wdat.worldOffset+234, Signed), nil
				}
				return -1, nil
			},
		},
	},
})
