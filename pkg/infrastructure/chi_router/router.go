package chi_router

import (
	"chi_boilerplate/pkg/adapters/db"
	"chi_boilerplate/pkg/adapters/repositories/sqlx_mysql"
	"chi_boilerplate/pkg/domain/services"
	"chi_boilerplate/pkg/domain/usecases"
	"chi_boilerplate/pkg/infrastructure/chi_router/handlers"
	"chi_boilerplate/pkg/infrastructure/chi_router/handlers/api"
	"chi_boilerplate/pkg/infrastructure/chi_router/handlers/web"
	"chi_boilerplate/pkg/infrastructure/logger"
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
	r := chi.NewRouter()

	// Middlewares
	s.initMiddlewares(r)

	// JWT token
	err := s.initJWTToken()
	if err != nil {
		return err
	}

	// Routes
	s.routes(r)

	fmt.Printf("Server started on %s:%s...\n", s.Addr, s.Port)
	return http.ListenAndServe(fmt.Sprintf("%s:%s", s.Addr, s.Port), r)
}

func (s *ChiServer) HandleError(f func(w http.ResponseWriter, r *http.Request) error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handlers.WrapError(f, s.Logger)(w, r)
	}
}

// routes defines the routes of the server.
func (s *ChiServer) routes(r *chi.Mux) {
	// Web routes
	r.Get("/health-check", s.HandleError(web.HealthCheck))

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
			// Public routes
			v1.Group(func(v1 chi.Router) {
				// User routes
				v1.Route("/", func(u chi.Router) {
					userRepo := sqlx_mysql.NewUserMysqlRepository(s.DB)
					userService := services.NewUser(userRepo)
					userUseCase := usecases.NewUser(userService)

					h := api.NewUser(u, s.Logger, userUseCase)
					h.UserPublicRoutes()
				})
			})

			// Protected routes
			v1.Group(func(v1 chi.Router) {
				s.initJWT(v1)

				// User routes
				v1.Route("/users", func(u chi.Router) {
					userRepo := sqlx_mysql.NewUserMysqlRepository(s.DB)
					userService := services.NewUser(userRepo)
					userUseCase := usecases.NewUser(userService)

					h := api.NewUser(u, s.Logger, userUseCase)
					h.UserProtectedRoutes()
				})
			})
		})
	})
}
