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
				wrldSection, err := SectionOffset("WRLD", 1)
				if err != nil {
					return nil, err
				}
				return worldData{worldOffset: wrldSection}, nil
			},
		},
		"worldSettings": &graphql.Field{
			Type:        worldDataType,
			Description: "Civ3 save data",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				wrldSection, err := SectionOffset("WRLD", 1)
				if err != nil {
					return nil, err
				}
				settings := WorldSettingsData{}
				settings.WorldSeed = ReadInt32(wrldSection+170, Signed)
				return settings, nil
			},
		},
	},
})
