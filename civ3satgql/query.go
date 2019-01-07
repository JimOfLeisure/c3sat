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
