package civ3satgql

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
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

// Handler wrapper to allow adding headers to all responses
// concept yoinked from http://echorand.me/dissecting-golangs-handlerfunc-handle-and-defaultservemux.html
func setHeaders(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set Origin headers for CORS
		// yoinked from http://stackoverflow.com/questions/12830095/setting-http-headers-in-golang Matt Bucci's answer
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers",
				"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		}
		// Since we're dynamically setting origin, don't let it get cached
		w.Header().Set("Vary", "Origin")
		handler.ServeHTTP(w, r)
	})
}

func Server(path string, bindAddress, bindPort string) error {
	var err error
	saveGame.data, _, err = parseciv3.ReadFile(path)
	if err != nil {
		return err
	}
	findSections()

	Schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
		// Mutation: MutationType,
	})
	if err != nil {
		return err
	}

	// create a graphl-go HTTP handler
	graphQlHandler := handler.New(&handler.Config{
		Schema: &Schema,
		Pretty: false,
		// GraphiQL provides simple web browser query interface pulled from Internet
		GraphiQL: false,
		// Playground provides fancier web browser query interface pulled from Internet
		Playground: true,
	})

	http.Handle("/graphql", setHeaders(graphQlHandler))
	log.Fatal(http.ListenAndServe(bindAddress+":"+bindPort, nil))
	return nil
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
