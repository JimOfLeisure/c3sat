package civ3satgql

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/myjimnelson/c3sat/civ3decompress"
	// "github.com/myjimnelson/c3sat/parseciv3"
)

type sectionType struct {
	name   string
	offset int
}

type saveGameType struct {
	path     string
	data     []byte
	sections []sectionType
}

var saveGame saveGameType
var defaultBic saveGameType
var currentBic saveGameType
var currentGame saveGameType

// populates the structure given a path to a sav file
func (sav *saveGameType) loadSave(path string) error {
	var err error
	sav.data, _, err = civ3decompress.ReadFile(path)
	if err != nil {
		return err
	}
	sav.path = path
	sav.populateSections()
	// complete hack, DELETEME
	currentBic = defaultBic
	gameOff, err := saveGame.sectionOffset("GAME", 2)
	if err != nil {
		return nil
	}
	currentGame.data = saveGame.data[gameOff-4:]
	currentGame.populateSections()
	// end complete hack
	return nil
}

// Find sections demarc'ed by 4-character ASCII headers and place into sections[]
func (sav *saveGameType) populateSections() {
	var i, count, offset int
	sav.sections = make([]sectionType, 0)
	// find sections demarc'ed by 4-character ASCII headers
	for i < len(sav.data) {
		if sav.data[i] < 0x20 || sav.data[i] > 0x5a {
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
			s.name = string(sav.data[offset:i])
			sav.sections = append(sav.sections, *s)
		}
	}
}

// returns just the filename part of the path assuming / or \ separators
func (sav *saveGameType) fileName() string {
	var o int
	for i := 0; i < len(sav.path); i++ {
		if sav.path[i] == 0x2f || sav.path[i] == 0x5c {
			o = i
		}
	}
	return sav.path[o+1:]
}

// Transitioning to this from the old SecionOffset stanalone function
func (sav *saveGameType) sectionOffset(sectionName string, nth int) (int, error) {
	var i, n int
	for i < len(sav.sections) {
		if sav.sections[i].name == sectionName {
			n++
			if n >= nth {
				return sav.sections[i].offset + len(sectionName), nil
			}
		}
		i++
	}
	return -1, errors.New("Could not find " + strconv.Itoa(nth) + " section named " + sectionName)
}

func (sav *saveGameType) readInt32(offset int, signed bool) int {
	n := int(sav.data[offset]) +
		int(sav.data[offset+1])*0x100 +
		int(sav.data[offset+2])*0x10000 +
		int(sav.data[offset+3])*0x1000000
	if signed && n&0x80000000 != 0 {
		n = -(n ^ 0xffffffff + 1)
	}
	return n
}

func (sav *saveGameType) readInt16(offset int, signed bool) int {
	n := int(sav.data[offset]) +
		int(sav.data[offset+1])*0x100
	if signed && n&0x8000 != 0 {
		n = -(n ^ 0xffff + 1)
	}
	return n
}

func (sav *saveGameType) readInt8(offset int, signed bool) int {
	n := int(sav.data[offset])
	if signed && n&0x80 != 0 {
		n = -(n ^ 0xff + 1)
	}
	return n
}

// ChangeSavePath updates the package saveGame structure with save file data at <path>
func ChangeSavePath(path string) error {
	err := saveGame.loadSave(path)
	return err
}

// ChangeSavePath updates the package saveGame structure with save file data at <path>
func ChangeDefaultBicPath(path string) error {
	err := defaultBic.loadSave(path)
	return err
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

// GraphQlHandler returns a GraphQL http handler
func GraphQlHandler() (http.Handler, error) {
	Schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
		// Mutation: MutationType,
	})
	if err != nil {
		return nil, err
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
	return graphQlHandler, nil
}

func Server(path string, bindAddress, bindPort string) error {
	err := saveGame.loadSave(path)
	if err != nil {
		return err
	}

	gQlHandler, err := GraphQlHandler()
	if err != nil {
		return err
	}

	http.Handle("/graphql", setHeaders(gQlHandler))
	http.ListenAndServe(bindAddress+":"+bindPort, nil)
	if err != nil {
		return err
	}
	return nil
}

func Query(query, path string) (string, error) {
	err := saveGame.loadSave(path)
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
	return string(out[:]), nil
}

// Adapting this from the parseciv3 module to pull data this module's way
//   To keep the existing seed command without having to gql everything
// WorldSettings Returns the information needed to regenerate the map
// presuming the map was originally generated by Civ3 and not later edited
func WorldSettings(path string) ([][3]string, error) {
	var worldSize = [...]string{"Tiny", "Small", "Standard", "Large", "Huge", "Random"}
	// // "No Barbarians" is actually -1. To make code simpler, add one to value for array index
	var barbs = [...]string{"No Barbarians", "Sedentary", "Roaming", "Restless", "Raging", "Random"}
	var landMass = [...]string{"Archipelago", "Continents", "Pangea", "Random"}
	var oceanCoverage = [...]string{"80% Water", "70% Water", "60% Water", "Random"}
	var climate = [...]string{"Arid", "Normal", "Wet", "Random"}
	var temperature = [...]string{"Warm", "Temperate", "Cool", "Random"}
	var age = [...]string{"3 Billion", "4 Billion", "5 Billion", "Random"}
	var settings [][3]string
	err := saveGame.loadSave(path)
	if err != nil {
		return nil, err
	}
	wrldSection, err := SectionOffset("WRLD", 1)
	if err != nil {
		return [][3]string{}, err
	}
	settings = [][3]string{
		{"Setting", "Choice", "Result"},
		{"World Seed", strconv.FormatUint(uint64(ReadInt32(wrldSection+170, Signed)), 10), ""},
		{"World Size", worldSize[ReadInt32(wrldSection+234, Signed)], ""},
		{"Barbarians", barbs[ReadInt32(wrldSection+194, Signed)+1], barbs[ReadInt32(wrldSection+198, Signed)+1]},
		{"Land Mass", landMass[ReadInt32(wrldSection+202, Signed)], landMass[ReadInt32(wrldSection+206, Signed)]},
		{"Water Coverage", oceanCoverage[ReadInt32(wrldSection+210, Signed)], oceanCoverage[ReadInt32(wrldSection+214, Signed)]},
		{"Climate", climate[ReadInt32(wrldSection+186, Signed)], climate[ReadInt32(wrldSection+190, Signed)]},
		{"Temperature", temperature[ReadInt32(wrldSection+218, Signed)], temperature[ReadInt32(wrldSection+222, Signed)]},
		{"Age", age[ReadInt32(wrldSection+226, Signed)], age[ReadInt32(wrldSection+230, Signed)]},
	}
	return settings, nil
}
