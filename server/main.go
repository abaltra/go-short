package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var ALPHABET = []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

var s *Store = NewStore()

func encode(id int) string {
	log.Printf("Encoding id %d", id)
	var sb strings.Builder
	l := len(ALPHABET)
	if id == 0 {
		sb.WriteRune(ALPHABET[0])
	}

	for id > 0 {
		sb.WriteRune(ALPHABET[id%l])
		id /= l
	}

	return reverse(sb.String())
}

func indexOf(c rune) int {
	for idx, r := range ALPHABET {
		if r == c {
			return idx
		}
	}

	return -1
}

func reverse(s string) string {
	r := []rune(s)

	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j+1 {
		r[i], r[j] = r[j], r[i]
	}

	return string(r)
}

func decode(code string) int {
	id := 0
	rcode := []rune(code)
	base := len(ALPHABET)
	for _, r := range rcode {
		id = id*base + indexOf(r)
	}

	return id
}

func GetForwardRoute(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	log.Printf("Trying to get entry %s from store\n", vars["routeID"])
	// TODO: get from persisten cache
	idx := decode(vars["routeID"])
	// url, ok := cache[vars["routeID"]]
	url := s.GetURL(idx)

	if url == "" {
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

type response struct {
	Code string
}

func ShortenRoute(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var p payload
	decoder.Decode(&p)
	if !strings.HasPrefix(p.URL, "http") {
		p.URL = fmt.Sprintf("https://%s", p.URL)
	}
	idx := s.InsertURL(p.URL)
	// TODO: Make this smarter. the hash should not be longer than 6 characters
	code := encode(idx)
	/// TODO: Move from in-memory to some persistent DB. Dynamo?
	// cache[code] = p.URL
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	log.Printf("Returning code %s for url %s\n", code, p.URL)
	resp := response{
		Code: code,
	}

	log.Print(resp)
	json.NewEncoder(w).Encode(resp)
	// fmt.Fprint(w, code)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", ShortenRoute).Methods("POST")
	router.HandleFunc("/{routeID}", GetForwardRoute).Methods("GET")

	c := cors.AllowAll()
	handler := c.Handler(router)

	srv := &http.Server{
		Handler: handler,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
