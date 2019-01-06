package civ3satgql

import (
	"encoding/json"

	"github.com/graphql-go/graphql"
	"github.com/myjimnelson/c3sat/parseciv3"
)

type saveGameType struct {
	data     []byte
	sections []string
}

var saveGame saveGameType

func Query(query, path string) (string, error) {
	var err error
	saveGame.data, _, err = parseciv3.ReadFile(path)
	if err != nil {
		return "", err
	}
	saveGame.sections = []string{"hello", "there"}
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
