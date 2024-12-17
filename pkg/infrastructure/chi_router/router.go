package chi_router

import (
	"chi_boilerplate/pkg/adapters/db"
	"chi_boilerplate/pkg/adapters/repositories/sqlx_mysql"
	"chi_boilerplate/pkg/domain/usecases"
	"chi_boilerplate/pkg/infrastructure/chi_router/handlers"
	"chi_boilerplate/pkg/infrastructure/chi_router/handlers/api"
	"chi_boilerplate/pkg/infrastructure/chi_router/handlers/web"
	"chi_boilerplate/pkg/infrastructure/logger"
	"chi_boilerplate/utils"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// ChiServer is a struct that represents a Chi server
type ChiServer struct {
	Addr   string
	Port   string
	DB     *db.MySQL
	Logger logger.CustomLogger
}

// NewChiServer creates a new ChiServer
func NewChiServer(addr, port string, db *db.MySQL, l logger.CustomLogger) ChiServer {
	return ChiServer{
		Addr:   addr,
		Port:   port,
		DB:     db,
		Logger: l,
	}
}

// Start the HTTP server
func (s *ChiServer) Start() error {
	r, err := s.Setup()
	if err != nil {
		return err
	}

	fmt.Printf("Server started on %s:%s...\n", s.Addr, s.Port)
	return http.ListenAndServe(fmt.Sprintf("%s:%s", s.Addr, s.Port), r)
}

// Setup the HTTP server
func (s *ChiServer) Setup() (*chi.Mux, error) {
	r := chi.NewRouter()

	// Middlewares
	s.initMiddlewares(r)

	// JWT token
	err := s.initJWTToken()
	if err != nil {
		return r, err
	}

	// Routes
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		utils.Err404(w, nil, "Ressource not found", nil)
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		utils.Err405(w, nil, "Method not allowed", nil)
	})
	s.routes(r)

	return r, nil
}

func (s *ChiServer) HandleError(f func(w http.ResponseWriter, r *http.Request) error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handlers.WrapError(f, s.Logger)(w, r)
	}
}

// routes defines the routes of the server.
func (s *ChiServer) routes(r *chi.Mux) {
	// Web routes
	r.Get("/health", s.HandleError(web.HealthCheck))

	// API documentation
	r.Route("/doc", func(d chi.Router) {
		d.Use(s.initBasicAuth())

		d.Get("/api-v1", s.HandleError(web.GetAPIv1Doc))
	})

	// Static files
	fs := http.FileServer(http.Dir("./assets"))
	r.Handle("/assets/*", http.StripPrefix("/assets/", fs))

	// API routes
	r.Route("/api", func(a chi.Router) {
		a.Use(s.initCORS())

		// Version 1
		a.Route("/v1", func(v1 chi.Router) {
			// User use case
			userRepo := sqlx_mysql.NewUserMysqlRepository(s.DB)
			userUseCase := usecases.NewUser(userRepo)

			// Public routes
			v1.Group(func(v1 chi.Router) {
				// User routes
				v1.Route("/", func(u chi.Router) {
					h := api.NewUser(u, s.Logger, userUseCase)
					h.UserPublicRoutes()
				})
			})

			// Protected routes
			v1.Group(func(v1 chi.Router) {
				s.initJWT(v1)

				// User routes
				v1.Route("/users", func(u chi.Router) {
					h := api.NewUser(u, s.Logger, userUseCase)
					h.UserProtectedRoutes()
				})
			})
		})
	})
}
