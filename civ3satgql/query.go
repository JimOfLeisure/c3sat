package civ3satgql

import (
	"github.com/graphql-go/graphql"
)

var queryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"hello": &graphql.Field{
			Type:        graphql.String,
			Description: "Get list of keys in LevelDB. Specifying multiple boolean arguments is invalid",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "world", nil
			},
		},
	},
})
