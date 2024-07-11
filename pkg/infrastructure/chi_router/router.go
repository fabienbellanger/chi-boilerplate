package chi_router

import (
	"chi_boilerplate/pkg/adapters/db"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// ChiServer is a struct that represents a Chi server
type ChiServer struct {
	Addr string
	Port string
	DB   *db.MySQL
}

// NewChiServer creates a new ChiServer
func NewChiServer(addr, port string, db *db.MySQL) ChiServer {
	return ChiServer{
		Addr: addr,
		Port: port,
		DB:   db,
	}
}

// Start the HTTP server
func (s ChiServer) Start() error {
	r := chi.NewRouter()

	// Middlewares
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

	// curl -H "Content-Type: application/json" -d '{"name":"xyz"}' --request POST "http://localhost:3000/hello"
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

	fmt.Printf("Server started on %s:%s...\n", s.Addr, s.Port)
	return http.ListenAndServe(fmt.Sprintf("%s:%s", s.Addr, s.Port), r)
}
