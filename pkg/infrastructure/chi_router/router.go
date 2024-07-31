package chi_router

import (
	"chi_boilerplate/pkg/adapters/db"
	"chi_boilerplate/pkg/adapters/repositories/sqlx_mysql"
	"chi_boilerplate/pkg/domain/services"
	"chi_boilerplate/pkg/domain/usecases"
	"chi_boilerplate/pkg/infrastructure/chi_router/handlers/api"
	"chi_boilerplate/pkg/infrastructure/chi_router/handlers/web"
	"chi_boilerplate/pkg/infrastructure/logger"
	"fmt"
	"go.uber.org/zap"
	"net/http"

	"github.com/go-chi/chi/v5"
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
func (s *ChiServer) Start() error {
	r := chi.NewRouter()

	// Custom logger
	zapLogger, err := logger.InitLogger()
	if err != nil {
		return err
	}

	// Middlewares
	s.initMiddlewares(r, zapLogger)

	// Routes
	s.routes(r, zapLogger)

	fmt.Printf("Server started on %s:%s...\n", s.Addr, s.Port)
	return http.ListenAndServe(fmt.Sprintf("%s:%s", s.Addr, s.Port), r)
}

// routes defines the routes of the server.
func (s *ChiServer) routes(r *chi.Mux, log *zap.Logger) {
	// Web routes
	r.Get("/health-check", web.HealthCheck)

	// API documentation
	r.Route("/doc", func(d chi.Router) {
		d.Use(initBasicAuth())

		d.Get("/api-v1", web.GetAPIv1Doc)
	})

	// Static files
	fs := http.FileServer(http.Dir("./assets"))
	r.Handle("/assets/*", http.StripPrefix("/assets/", fs))

	// API routes
	r.Route("/api", func(a chi.Router) {
		a.Use(initCORS())

		// Version 1
		a.Route("/v1", func(v1 chi.Router) {
			// Public routes
			v1.Group(func(v1 chi.Router) {
				// User routes
				v1.Route("/users", func(u chi.Router) {
					userRepo := sqlx_mysql.NewUserMysqlRepository(s.DB)
					userService := services.NewUser(userRepo)
					userUseCase := usecases.NewUser(userService)

					h := api.NewUser(u, userUseCase)
					h.UserProtectedRoutes()
				})
			})

			// Protected routes
			v1.Group(func(v1 chi.Router) {
				// User routes
				v1.Route("/", func(u chi.Router) {
					userRepo := sqlx_mysql.NewUserMysqlRepository(s.DB)
					userService := services.NewUser(userRepo)
					userUseCase := usecases.NewUser(userService)

					h := api.NewUser(u, userUseCase)
					h.UserPublicRoutes()
				})
			})
		})
	})
}
