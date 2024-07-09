package services

import (
	"chi_boilerplate/pkg/domain/entities"
	"chi_boilerplate/pkg/domain/repositories"
	"chi_boilerplate/pkg/domain/requests"
	"chi_boilerplate/pkg/domain/responses"
	values_objects "chi_boilerplate/pkg/domain/value_objects"
	"chi_boilerplate/utils"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

// UserService is an interface for user service
type UserService interface {
	Login(req requests.UserLogin) (responses.UserLogin, *utils.HTTPError)
	Create(req requests.UserCreation) (entities.User, *utils.HTTPError)
	// GetAll(req requests.Pagination) (responses.UsersListPaginated, *utils.HTTPError)
	// GetByID(id requests.UserByID) (entities.User, *utils.HTTPError)
	// Delete(id requests.UserByID) *utils.HTTPError
	// Update(req requests.UserUpdate) (entities.User, *utils.HTTPError)
}

type userService struct {
	userRepository repositories.UserRepository
}

// NewUser returns a new user service
func NewUser(repo repositories.UserRepository) UserService {
	return &userService{repo}
}

// Login user
func (us userService) Login(req requests.UserLogin) (responses.UserLogin, *utils.HTTPError) {
	loginErrors := utils.ValidateStruct(req)
	if loginErrors != nil {
		return responses.UserLogin{}, utils.NewHTTPError(utils.StatusBadRequest, "Invalid body", loginErrors, nil)
	}

	user, err := us.userRepository.Login(req.Email, req.Password)
	if err != nil {
		var e *utils.HTTPError
		if errors.Is(err, repositories.ErrUserNotFound) {
			e = utils.NewHTTPError(utils.StatusUnauthorized, "Unauthorized", nil, nil)
		} else {
			e = utils.NewHTTPError(utils.StatusInternalServerError, "Internal server error", "Error during authentication", err)
		}
		return responses.UserLogin{}, e
	}

	// Create token
	token, expiresAt, err := user.GenerateJWT(
		viper.GetDuration("JWT_LIFETIME"),
		viper.GetString("JWT_ALGO"),
		viper.GetString("JWT_SECRET"))
	if err != nil {
		return responses.UserLogin{}, utils.NewHTTPError(utils.StatusInternalServerError, "Internal server error", "Error during token generation", err)
	}

	return responses.UserLogin{
		User:      user,
		Token:     token,
		ExpiresAt: expiresAt.Format(time.RFC3339),
	}, nil
}

// Create user
func (us userService) Create(req requests.UserCreation) (entities.User, *utils.HTTPError) {
	creationErrors := utils.ValidateStruct(req)
	if creationErrors != nil {
		return entities.User{}, utils.NewHTTPError(utils.StatusBadRequest, "Invalid body", creationErrors, nil)
	}

	email, err := values_objects.NewEmail(req.Email)
	if err != nil {
		return entities.User{}, utils.NewHTTPError(utils.StatusBadRequest, "Invalid email", "Error during user creation", err)
	}
	password, err := values_objects.NewPassword(req.Password)
	if err != nil {
		return entities.User{}, utils.NewHTTPError(utils.StatusBadRequest, "Invalid password", "Error during user creation", err)
	}

	user := entities.User{
		ID:        uuid.New(),
		Lastname:  req.Lastname,
		Firstname: req.Firstname,
		Password:  password,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := us.userRepository.Create(&user); err != nil {
		return entities.User{}, utils.NewHTTPError(utils.StatusInternalServerError, "Database error", "Error during user creation", err)
	}

	return user, nil
}
