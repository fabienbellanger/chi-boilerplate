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
	Create(req requests.UserCreation) (responses.UserCreation, *utils.HTTPError)
	GetByID(id requests.UserByID) (responses.UserById, *utils.HTTPError)
	// GetAll(req requests.Pagination) (responses.UsersListPaginated, *utils.HTTPError)
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
func (us *userService) Login(req requests.UserLogin) (responses.UserLogin, *utils.HTTPError) {
	loginErrors := utils.ValidateStruct(req)
	if loginErrors != nil {
		return responses.UserLogin{}, utils.NewHTTPError(utils.StatusBadRequest, "Invalid request data", loginErrors, nil)
	}

	hashedPassword := entities.HashUserPassword(req.Password)

	userRepo, err := us.userRepository.Login(requests.UserLogin{Email: req.Email, Password: hashedPassword})
	if err != nil {
		var e *utils.HTTPError
		if errors.Is(err, repositories.ErrUserNotFound) {
			e = utils.NewHTTPError(utils.StatusUnauthorized, "Unauthorized", nil, nil)
		} else {
			e = utils.NewHTTPError(utils.StatusInternalServerError, "Internal server error", "Error during authentication", err)
		}
		return responses.UserLogin{}, e
	}

	user, err := userRepo.ToUser()
	if err != nil {
		return responses.UserLogin{}, utils.NewHTTPError(utils.StatusInternalServerError, "Internal server error", "Error during authentication", err)
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
		ID:        user.ID,
		Email:     user.Email.Value,
		Lastname:  user.Lastname,
		Firstname: user.Firstname,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		Token:     token,
		ExpiresAt: expiresAt.Format(time.RFC3339),
	}, nil
}

// Create user
func (us *userService) Create(req requests.UserCreation) (responses.UserCreation, *utils.HTTPError) {
	creationErrors := utils.ValidateStruct(req)
	if creationErrors != nil {
		return responses.UserCreation{}, utils.NewHTTPError(utils.StatusBadRequest, "Invalid request data", creationErrors, nil)
	}

	now := time.Now()
	userID := uuid.New()
	user := requests.UserCreationRepository{
		ID:        userID.String(),
		Lastname:  req.Lastname,
		Firstname: req.Firstname,
		Password:  entities.HashUserPassword(req.Password),
		Email:     req.Email,
		CreatedAt: now.Format(utils.SqlDateTimeFormat),
		UpdatedAt: now.Format(utils.SqlDateTimeFormat),
	}

	if err := us.userRepository.Create(user); err != nil {
		return responses.UserCreation{}, utils.NewHTTPError(utils.StatusInternalServerError, "Database error", "Error during user creation", err)
	}
	email, err := values_objects.NewEmail(user.Email)
	if err != nil {
		return responses.UserCreation{}, utils.NewHTTPError(utils.StatusInternalServerError, "Internal server error", "Error when getting user ", err)
	}

	return responses.UserCreation{
		ID:        userID,
		Email:     email,
		Lastname:  user.Lastname,
		Firstname: user.Firstname,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// GetByID user
func (us *userService) GetByID(req requests.UserByID) (responses.UserById, *utils.HTTPError) {
	reqErrors := utils.ValidateStruct(req)
	if reqErrors != nil {
		return responses.UserById{}, utils.NewHTTPError(utils.StatusBadRequest, "Invalid request data", reqErrors, nil)
	}

	userRepo, err := us.userRepository.GetByID(requests.UserByID{ID: req.ID})
	if err != nil {
		var e *utils.HTTPError
		if errors.Is(err, repositories.ErrUserNotFound) {
			e = utils.NewHTTPError(utils.StatusNotFound, "User not found", nil, nil)
		} else {
			e = utils.NewHTTPError(utils.StatusInternalServerError, "Internal server error", "Error when getting user ", err)
		}
		return responses.UserById{}, e
	}

	email, err := values_objects.NewEmail(userRepo.Email)
	if err != nil {
		return responses.UserById{}, utils.NewHTTPError(utils.StatusInternalServerError, "Internal server error", "Error when getting user ", err)
	}
	createdAt, err := time.Parse(time.RFC3339, userRepo.CreatedAt)
	if err != nil {
		return responses.UserById{}, utils.NewHTTPError(utils.StatusInternalServerError, "Internal server error", "Error when getting user ", err)
	}
	updatedAt, err := time.Parse(time.RFC3339, userRepo.UpdatedAt)
	if err != nil {
		return responses.UserById{}, utils.NewHTTPError(utils.StatusInternalServerError, "Internal server error", "Error when getting user ", err)
	}
	user := responses.UserById{
		ID:        uuid.MustParse(userRepo.ID),
		Email:     email,
		Lastname:  userRepo.Lastname,
		Firstname: userRepo.Firstname,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return user, nil
}
