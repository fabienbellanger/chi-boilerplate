package chi_router

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func Start() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	// curl "http://localhost:3000/hello/fabien?q=test"
	r.Get("/hello/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		query := r.URL.Query().Get("q")

		w.Write([]byte("Hello " + name + " with " + query + "!"))
	})

	// curl -H "Content-Type: application/json" -d '{"name":"xyz"}' --request POST http://localhost:3000/hello
	r.Post("/hello", func(w http.ResponseWriter, r *http.Request) {
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
	})

	http.ListenAndServe("localhost:3000", r)
}
