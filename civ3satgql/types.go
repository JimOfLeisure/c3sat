package civ3satgql

import (
	"github.com/graphql-go/graphql"
)

var civ3Type = graphql.NewObject(graphql.ObjectConfig{
	Name: "civ3",
	Fields: graphql.Fields{
		"worldSeed": &graphql.Field{
			Type:        graphql.Int,
			Description: "Random seed of rando worlds",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var i int
				for i < len(saveGame.sections) {
					if saveGame.sections[i].name == "WRLD" {
						myOffset := saveGame.sections[i].offset + 174
						foo := readInt32(myOffset, signed)
						return foo, nil
					}
					i++
				}
				return 0, nil
			},
		},
	},
})
