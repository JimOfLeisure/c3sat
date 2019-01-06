package civ3satgql

import (
	"github.com/graphql-go/graphql"
)

func readInt32(offset int) int {
	return int(saveGame.data[offset]) +
		int(saveGame.data[offset+1])*0x100 +
		int(saveGame.data[offset+2])*0x10000 +
		int(saveGame.data[offset+3])*0x1000000

}

var queryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"civ3": &graphql.Field{
			Type:        graphql.Int,
			Description: "Testing",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var i int
				for i < len(saveGame.sections) {
					if saveGame.sections[i].name == "WRLD" {
						myOffset := saveGame.sections[i].offset + 174
						foo := readInt32(myOffset)
						return foo, nil
					}
					i++
				}
				return 0, nil
			},
		},
	},
})
