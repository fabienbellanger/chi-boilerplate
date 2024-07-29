package api

import (
	"chi_boilerplate/pkg/domain/requests"
	"chi_boilerplate/pkg/domain/usecases"
	"chi_boilerplate/utils"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// User handler
type User struct {
	router      chi.Router
	userUseCase usecases.User
}

// NewUser returns a new Handler
func NewUser(r chi.Router, userUseCase usecases.User) User {
	return User{
		router:      r,
		userUseCase: userUseCase,
	}
}

// UserPublicRoutes adds users public routes
func (u *User) UserPublicRoutes() {
	u.router.Post("/login", u.login)
}

// UserProtectedRoutes adds users protected routes
func (u *User) UserProtectedRoutes() {
	u.router.Post("/", u.create)
}

func (u *User) login(w http.ResponseWriter, r *http.Request) {
	var req requests.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Err400(w, err, "Error decoding body")
		return
	}

	res, err := u.userUseCase.Login(req)
	if err != nil {
		err.SendError(w)
		return
	}

	utils.JSON(w, res)
}

func (u *User) create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Register"))
}
