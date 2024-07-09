package usecases

import (
	"chi_boilerplate/pkg/domain/entities"
	"chi_boilerplate/pkg/domain/requests"
	"chi_boilerplate/pkg/domain/responses"
	"chi_boilerplate/pkg/domain/services"
	"chi_boilerplate/utils"
)

// User is an interface for user use cases
type User interface {
	Login(req requests.UserLogin) (responses.UserLogin, *utils.HTTPError)
	Create(req requests.UserCreation) (entities.User, *utils.HTTPError)
	// GetAll(req requests.Pagination) (responses.UsersListPaginated, *utils.HTTPError)
	// GetByID(id requests.UserByID) (entities.User, *utils.HTTPError)
	// Delete(id requests.UserByID) *utils.HTTPError
	// Update(req requests.UserUpdate) (entities.User, *utils.HTTPError)
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
func (uc *userUseCase) Create(req requests.UserCreation) (entities.User, *utils.HTTPError) {
	return uc.userService.Create(req)
}

// // GetAll users
// func (uc *userUseCase) GetAll(req requests.Pagination) (responses.UsersListPaginated, *utils.HTTPError) {
// 	return uc.userService.GetAll(req)
// }

// // GetByID user
// func (uc *userUseCase) GetByID(id requests.UserByID) (entities.User, *utils.HTTPError) {
// 	return uc.userService.GetByID(id)
// }

// // Delete user
// func (uc *userUseCase) Delete(id requests.UserByID) *utils.HTTPError {
// 	return uc.userService.Delete(id)
// }

// // Update user
// func (uc *userUseCase) Update(req requests.UserUpdate) (entities.User, *utils.HTTPError) {
// 	return uc.userService.Update(req)
// }
