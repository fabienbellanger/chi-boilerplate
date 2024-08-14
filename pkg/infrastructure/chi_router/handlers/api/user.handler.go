package api

import (
	"chi_boilerplate/pkg/domain/requests"
	"chi_boilerplate/pkg/domain/usecases"
	"chi_boilerplate/pkg/infrastructure/chi_router/handlers"
	"chi_boilerplate/pkg/infrastructure/logger"
	"chi_boilerplate/utils"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// User handler
type User struct {
	router      chi.Router
	userUseCase usecases.User
	logger      logger.CustomLogger
}

// NewUser returns a new Handler
func NewUser(r chi.Router, l logger.CustomLogger, userUseCase usecases.User) User {
	return User{
		router:      r,
		userUseCase: userUseCase,
		logger:      l,
	}
}

// UserPublicRoutes adds users public routes
func (u *User) UserPublicRoutes() {
	u.router.Post("/login", handlers.WrapError(u.login, u.logger))
}

// UserProtectedRoutes adds users protected routes
func (u *User) UserProtectedRoutes() {
	u.router.Post("/", handlers.WrapError(u.create, u.logger))
	u.router.Get("/{id}", handlers.WrapError(u.getByID, u.logger))
	u.router.Delete("/{id}", handlers.WrapError(u.delete, u.logger))
}

func (u *User) login(w http.ResponseWriter, r *http.Request) error {
	var body requests.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return utils.Err400(w, err, "Error decoding body", nil)
	}

	res, err := u.userUseCase.Login(body)
	if err != nil {
		return err.SendError(w)
	}

	return utils.JSON(w, res)
}

func (u *User) create(w http.ResponseWriter, r *http.Request) error {
	var body requests.UserCreation
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return utils.Err400(w, err, "Error decoding body", nil)
	}

	res, err := u.userUseCase.Create(body)
	if err != nil {
		return err.SendError(w)
	}

	return utils.JSON(w, res.ToUserHTTP())
}

func (u *User) getByID(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	if id == "" {
		return utils.Err400(w, nil, "ID is required", nil)
	}

	res, err := u.userUseCase.GetByID(requests.UserByID{ID: id})
	if err != nil {
		return err.SendError(w)
	}

	return utils.JSON(w, res.ToUserHTTP())
}

func (u *User) delete(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	if id == "" {
		return utils.Err400(w, nil, "ID is required", nil)
	}

	err := u.userUseCase.Delete(requests.UserDelete{ID: id})
	if err != nil {
		return err.SendError(w)
	}

	return utils.NoContent(w)
}
