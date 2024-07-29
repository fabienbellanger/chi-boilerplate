package web

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// HealthCheck returns status code 200.
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// TODO: Remove this handler
func Panic(w http.ResponseWriter, r *http.Request) {
	panic("Panic test")
}

// TODO: Remove this handler
func GetHello(w http.ResponseWriter, r *http.Request) {
	// curl "http://localhost:3000/hello/fabien?q=test"
	name := chi.URLParam(r, "name")
	query := r.URL.Query().Get("q")

	w.Write([]byte("Hello " + name + " with q=" + query))
}

// TODO: Remove this handler
func PostHello(w http.ResponseWriter, r *http.Request) {
	// curl -H "Content-Type: application/json" -d '{"name":"xyz"}' --request POST "http://localhost:3000/hello"
	type Person struct {
		Name string
	}

	var p Person
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error decoding body"))
		return
	}

	w.Write([]byte("Hello " + p.Name + "!"))
}
