package usecases

import (
	"chi_boilerplate/pkg/domain/requests"
	"chi_boilerplate/pkg/domain/responses"
	"chi_boilerplate/pkg/domain/services"
	"chi_boilerplate/utils"
)

// User is an interface for user use cases
type User interface {
	Login(requests.UserLogin) (responses.UserLogin, *utils.HTTPError)
	Create(requests.UserCreation) (responses.UserCreation, *utils.HTTPError)
	GetByID(requests.UserByID) (responses.UserById, *utils.HTTPError)
	GetAll(requests.UsersList) (responses.UsersList, *utils.HTTPError)
	Delete(requests.UserDelete) *utils.HTTPError
	// Update(requests.UserUpdate) (entities.User, *utils.HTTPError)
}

type userUseCase struct {
	userService services.UserService
}

// NewUser returns a new User use case
func NewUser(userService services.UserService) User {
	return &userUseCase{userService}
}

// Login user
func (uc *userUseCase) Login(req requests.UserLogin) (responses.UserLogin, *utils.HTTPError) {
	return uc.userService.Login(req)
}

// Create user
func (uc *userUseCase) Create(req requests.UserCreation) (responses.UserCreation, *utils.HTTPError) {
	return uc.userService.Create(req)
}

// GetByID user
func (uc *userUseCase) GetByID(id requests.UserByID) (responses.UserById, *utils.HTTPError) {
	return uc.userService.GetByID(id)
}

// GetAll users
func (uc *userUseCase) GetAll(req requests.UsersList) (responses.UsersList, *utils.HTTPError) {
	return uc.userService.GetAll(req)
}

// Delete user
func (uc *userUseCase) Delete(req requests.UserDelete) *utils.HTTPError {
	return uc.userService.Delete(req)
}

// // Update user
// func (uc *userUseCase) Update(req requests.UserUpdate) (entities.User, *utils.HTTPError) {
// 	return uc.userService.Update(req)
// }
