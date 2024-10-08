package usecases

import (
	"chi_boilerplate/pkg/domain/entities"
	"chi_boilerplate/pkg/domain/repositories"
	"chi_boilerplate/pkg/domain/requests"
	"chi_boilerplate/pkg/domain/responses"
	vo "chi_boilerplate/pkg/domain/value_objects"
	"chi_boilerplate/utils"
	"errors"
	"time"

	"github.com/spf13/viper"
)

// User is an interface for user use cases
type User interface {
	Login(requests.UserLogin) (responses.UserLogin, *utils.HTTPError)
	Create(requests.UserCreation) (responses.UserCreation, *utils.HTTPError)
	GetByID(requests.UserByID) (responses.UserById, *utils.HTTPError)
	GetAll(requests.UsersList) (responses.UsersList, *utils.HTTPError)
	Delete(requests.UserDelete) *utils.HTTPError
	Update(requests.UserUpdate) (responses.UserById, *utils.HTTPError)
}

type userUseCase struct {
	userRepository repositories.UserRepository
}

// NewUser returns a new User use case
func NewUser(userRepository repositories.UserRepository) User {
	return &userUseCase{userRepository}
}

// Login user
func (uc *userUseCase) Login(req requests.UserLogin) (responses.UserLogin, *utils.HTTPError) {
	loginErrors := utils.ValidateStruct(req)
	if loginErrors != nil {
		return responses.UserLogin{}, utils.NewHTTPError(utils.StatusBadRequest, "Invalid request data", loginErrors, nil)
	}

	hashedPassword := entities.HashUserPassword(req.Password)

	userRepo, err := uc.userRepository.Login(requests.UserLogin{Email: req.Email, Password: hashedPassword})
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
		AccessToken:          token,
		AccessTokenExpiresAt: expiresAt.Format(time.RFC3339),
	}, nil
}

// Create user
func (uc *userUseCase) Create(req requests.UserCreation) (responses.UserCreation, *utils.HTTPError) {
	creationErrors := utils.ValidateStruct(req)
	if creationErrors != nil {
		return responses.UserCreation{}, utils.NewHTTPError(utils.StatusBadRequest, "Invalid request data", creationErrors, nil)
	}

	now := time.Now()
	userID := vo.NewID()
	user := requests.UserCreationRepository{
		ID:        userID.String(),
		Lastname:  req.Lastname,
		Firstname: req.Firstname,
		Password:  entities.HashUserPassword(req.Password),
		Email:     req.Email,
		CreatedAt: now.Format(utils.SqlDateTimeFormat),
		UpdatedAt: now.Format(utils.SqlDateTimeFormat),
	}

	if err := uc.userRepository.Create(user); err != nil {
		return responses.UserCreation{}, utils.NewHTTPError(utils.StatusInternalServerError, "Database error", "Error during user creation", err)
	}
	email, err := vo.NewEmail(user.Email)
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
func (uc *userUseCase) GetByID(req requests.UserByID) (responses.UserById, *utils.HTTPError) {
	reqErrors := utils.ValidateStruct(req)
	if reqErrors != nil {
		return responses.UserById{}, utils.NewHTTPError(utils.StatusBadRequest, "Invalid request data", reqErrors, nil)
	}

	id, err := vo.NewIDFrom(req.ID)
	if err != nil {
		return responses.UserById{}, utils.NewHTTPError(utils.StatusBadRequest, "Invalid request data", reqErrors, nil)
	}

	userRepo, err := uc.userRepository.GetByID(requests.UserByID{ID: req.ID})
	if err != nil {
		var e *utils.HTTPError
		if errors.Is(err, repositories.ErrUserNotFound) {
			e = utils.NewHTTPError(utils.StatusNotFound, "User not found", nil, nil)
		} else {
			e = utils.NewHTTPError(utils.StatusInternalServerError, "Internal server error", "Error when getting user ", err)
		}
		return responses.UserById{}, e
	}

	email, err := vo.NewEmail(userRepo.Email)
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
		ID:        id,
		Email:     email,
		Lastname:  userRepo.Lastname,
		Firstname: userRepo.Firstname,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return user, nil
}

// Delete user
func (uc *userUseCase) Delete(req requests.UserDelete) *utils.HTTPError {
	reqErrors := utils.ValidateStruct(req)
	if reqErrors != nil {
		return utils.NewHTTPError(utils.StatusBadRequest, "Invalid request data", reqErrors, nil)
	}

	err := uc.userRepository.Delete(requests.UserDelete{ID: req.ID})
	if err != nil {
		var e *utils.HTTPError
		if errors.Is(err, repositories.ErrUserNotFound) {
			e = utils.NewHTTPError(utils.StatusNotFound, "User not found", nil, nil)
		} else {
			e = utils.NewHTTPError(utils.StatusInternalServerError, "Internal server error", "Error when getting user ", err)
		}
		return e
	}

	return nil
}

// GetAll returns all users with pagination
func (uc *userUseCase) GetAll(req requests.UsersList) (responses.UsersList, *utils.HTTPError) {
	var list responses.UsersList
	users, err := uc.userRepository.GetAll(req)
	if err != nil {
		return responses.UsersList{}, utils.NewHTTPError(utils.StatusInternalServerError, "Internal server error", "Error when getting all users", err)
	}
	list.Data = users

	total, err := uc.userRepository.CountAll()
	if err != nil {
		return responses.UsersList{}, utils.NewHTTPError(utils.StatusInternalServerError, "Internal server error", "Error when getting all users", err)
	}
	list.Total = total

	return list, nil
}

// Update user
func (uc *userUseCase) Update(req requests.UserUpdate) (responses.UserById, *utils.HTTPError) {
	reqErrors := utils.ValidateStruct(req)
	if reqErrors != nil {
		return responses.UserById{}, utils.NewHTTPError(utils.StatusBadRequest, "Invalid request data", reqErrors, nil)
	}

	err := uc.userRepository.Update(requests.UserUpdateRepository{
		ID:        req.ID,
		Lastname:  req.Lastname,
		Firstname: req.Firstname,
		Email:     req.Email,
		Password:  entities.HashUserPassword(req.Password),
		UpdatedAt: time.Now().Format(utils.SqlDateTimeFormat),
	})
	if err != nil {
		return responses.UserById{}, utils.NewHTTPError(utils.StatusInternalServerError, "Error when updating user", err, nil)
	}

	return uc.GetByID(requests.UserByID{ID: req.ID})
}
