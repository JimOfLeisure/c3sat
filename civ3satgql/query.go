package civ3satgql

import (
	"github.com/graphql-go/graphql"
)

var queryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"civ3": &graphql.Field{
			Type:        civ3Type,
			Description: "Civ3 save data",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// I have to return something else subfields won't return
				return "test", nil
			},
		},
	},
})
