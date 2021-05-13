package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var cache map[string]string

func GetForwardRoute(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	log.Printf("Trying to get entry %s from cache %v\n", vars["routeID"], cache)
	// TODO: get from persisten cache
	url, ok := cache[vars["routeID"]]

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// TODO: Keep track of metrics? Delay this or do we do it in the same loop? goroutine?
	log.Printf("Redirecting %s to %s!\n", vars["routeID"], url)
	http.Redirect(w, r, url, http.StatusPermanentRedirect)
}

type payload struct {
	URL string
}

func ShortenRoute(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var p payload
	decoder.Decode(&p)
	// TODO: Make this smarter. the hash should not be longer than 6 characters
	code := uuid.NewString()
	/// TODO: Move from in-memory to some persistent DB. Dynamo?
	cache[code] = p.URL
	w.WriteHeader(http.StatusOK)
	log.Printf("Returning code %s for url %s\n", code, p.URL)
	fmt.Fprint(w, code)
}

func main() {
	cache = make(map[string]string)
	router := mux.NewRouter()
	router.HandleFunc("/", ShortenRoute).Methods("POST")
	router.HandleFunc("/{routeID}", GetForwardRoute).Methods("GET")

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
