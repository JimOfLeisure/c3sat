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
				return 0, nil
			},
		},
	},
})
