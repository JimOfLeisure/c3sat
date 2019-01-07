package civ3satgql

import (
	"github.com/graphql-go/graphql"
)

type worldData struct {
	worldOffset int
}

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
