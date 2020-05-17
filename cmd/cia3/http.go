package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/markbates/pkger"
	"github.com/myjimnelson/c3sat/queryciv3"
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

// Set Content-Type explicitly since net/http FileServer seems to use Win registry which on many systems has wrong info for js and css in particular
func setContentTypeHeaders(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := r.RequestURI
		if len(s) > 3 && strings.ToLower(s[len(s)-3:]) == ".js" {
			w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		} else if len(s) > 4 && strings.ToLower(s[len(s)-4:]) == ".css" {
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
		} else if len(s) > 5 && strings.ToLower(s[len(s)-5:]) == ".html" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		}
		handler.ServeHTTP(w, r)
	})
}

func server() {
	gQlHandler, err := queryciv3.GraphQlHandler()
	if err != nil {
		log.Fatal(err)
	}
	staticFiles := http.FileServer(pkger.Dir("/cmd/cia3/html"))
	http.Handle("/", setContentTypeHeaders(staticFiles))
	http.Handle("/graphql", setHeaders(gQlHandler))
	http.Handle("/events", setHeaders(http.Handler(http.HandlerFunc(longPoll.SubscriptionHandler))))
	err = http.ListenAndServe(addr+":"+httpPort, nil)
	if err != nil {
		log.Fatal(err)
	}
}
