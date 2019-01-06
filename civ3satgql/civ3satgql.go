package civ3satgql

import (
	"encoding/json"

	"github.com/graphql-go/graphql"
	"github.com/myjimnelson/c3sat/parseciv3"
)

type sectionType struct {
	name   string
	offset int
	length int
}

type saveGameType struct {
	data     []byte
	sections []sectionType
}

var saveGame saveGameType

func findSections() {
	var i, count, offset int
	for i < len(saveGame.data) {
		// for i < 83000 {
		if saveGame.data[i] < 0x20 || saveGame.data[i] > 0x5a {
			count = 0
		} else {
			if count == 0 {
				offset = i
			}
			count++
		}
		i++
		if count > 3 {
			count = 0
			s := new(sectionType)
			s.offset = offset
			s.name = string(saveGame.data[offset:i])
			saveGame.sections = append(saveGame.sections, *s)
			// fmt.Println(string(saveGame.data[offset:i]) + " " + strconv.Itoa(offset))
		}
	}
}

func Query(query, path string) (string, error) {
	var err error
	saveGame.data, _, err = parseciv3.ReadFile(path)
	if err != nil {
		return "", err
	}
	findSections()
	// fmt.Println(saveGame.sections[len(saveGame.sections)-1])
	// saveGame.sections = []string{"hello", "there"}
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
