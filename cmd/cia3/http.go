package main

import (
	"log"
	"net/http"

	"github.com/myjimnelson/c3sat/civ3satgql"
)

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
	http.Handle("/graphql", setHeaders(gQlHandler))
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
