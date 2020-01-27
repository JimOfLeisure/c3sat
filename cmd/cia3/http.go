package main

import (
	"log"
	"net/http"

	"github.com/markbates/pkger"
	"github.com/myjimnelson/c3sat/civ3satgql"
)

const addr = "127.0.0.1"

var httpUrlString string
var httpPort = "8080"
var httpPortTry = []string{
	":8080",
	":8000",
	":8888",
	":0",
}

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

func server() {
	gQlHandler, err := civ3satgql.GraphQlHandler()
	if err != nil {
		log.Fatal(err)
	}
	staticFiles := http.FileServer(pkger.Dir("/cmd/cia3/html"))
	http.Handle("/", staticFiles)
	http.Handle("/graphql", setHeaders(gQlHandler))
	http.Handle("/events", setHeaders(http.Handler(http.HandlerFunc(longPoll.SubscriptionHandler))))
	err = http.ListenAndServe(addr+":"+httpPort, nil)
	if err != nil {
		log.Fatal(err)
	}
}
