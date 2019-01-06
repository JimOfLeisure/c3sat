package civ3satgql

import (
	"github.com/graphql-go/graphql"
)

var civ3Type = graphql.NewObject(graphql.ObjectConfig{
	Name: "civ3",
	Fields: graphql.Fields{
		"worldSeed": &graphql.Field{
			Type:        graphql.Int,
			Description: "Random seed of random worlds",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				mySection, err := sectionOffset("WRLD", 1)
				if err != nil {
					return 0, nil
				}
				return readInt32(mySection+170, signed), nil
			},
		},
	},
})
