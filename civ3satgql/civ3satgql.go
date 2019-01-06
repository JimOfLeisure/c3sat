package civ3satgql

import (
	"encoding/json"

	"github.com/graphql-go/graphql"
	"github.com/myjimnelson/c3sat/parseciv3"
)

var saveGame []byte

func Query(query, path string) (string, error) {
	var err error
	saveGame, _, err = parseciv3.ReadFile(path)
	if err != nil {
		return "", err
	}
	Schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
		// Mutation: MutationType,
	})
	if err != nil {
		return "", err
	}
	result := graphql.Do(graphql.Params{
		Schema:        Schema,
		RequestString: query,
	})
	out, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	// return hex.EncodeToString(saveGame[:4]), nil
	return string(out[:]), nil
}
