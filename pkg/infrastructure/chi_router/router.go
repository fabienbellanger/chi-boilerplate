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
	logger, err := logger.InitLogger()
	if err != nil {
		return err
	}

	// Middlewares
	s.initMiddlewares(r, logger)

	// Web routes
	r.Get("/health-check", web.HealthCheck)
	r.Get("/panic", web.Panic)
	r.Get("/hello/{name}", web.GetHello)
	r.Post("/hello", web.PostHello)

	// API routes
	r.Route("/api/v1", func(v1 chi.Router) {
		// User routes
		v1.Route("/users", func(u chi.Router) {
			userRepo := sqlx_mysql.NewUserMysqlRepository(s.DB)
			userService := services.NewUser(userRepo)
			userUseCase := usecases.NewUser(userService)

			h := api.NewUser(u, userUseCase)
			h.UserPublicRoutes()
			h.UserProtectedRoutes()
		})
	}).With(initCORS())

	fmt.Printf("Server started on %s:%s...\n", s.Addr, s.Port)
	return http.ListenAndServe(fmt.Sprintf("%s:%s", s.Addr, s.Port), r)
}
