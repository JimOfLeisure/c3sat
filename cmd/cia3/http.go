package main

import (
	"log"
	"net/http"

	"github.com/markbates/pkger"
	"github.com/myjimnelson/c3sat/civ3satgql"
)

const addr = "127.0.0.1"
const port = "8080"

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
	staticFiles := http.FileServer(pkger.Dir("github.com/myjimnelson/c3sat:/cmd/cia3/html"))
	// Can't figure out how to make pkger work for non-root
	http.Handle("/", staticFiles)
	http.Handle("/graphql", setHeaders(gQlHandler))
	http.Handle("/events", setHeaders(http.Handler(http.HandlerFunc(longPoll.SubscriptionHandler))))
	// fmt.Println("Opening local web server, please browse to http://" + addr + ":" + port + "/isocss.html")
	// fmt.Println("Press control-C in this window or close it to end program")
	err = http.ListenAndServe(addr+":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
