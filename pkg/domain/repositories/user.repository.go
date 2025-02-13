package repositories

import (
	"chi_boilerplate/pkg/domain/requests"
	"chi_boilerplate/pkg/domain/responses"
	"errors"
)

var (
	// ErrUserNotFound is the error returned when a user is not found.
	ErrUserNotFound = errors.New("user not found")
)

// UserRepository is the interface that wraps the basic user repository methods.
type UserRepository interface {
	Create(requests.UserCreationRepository) error
	GetByID(requests.UserByID) (responses.UserByIdRepository, error)
	GetAll(requests.UsersList) ([]responses.UsersListRepository, error)
	CountAll() (int64, error)
	Delete(requests.UserDelete) error
	Update(requests.UserUpdateRepository) error
	GetByEmail(requests.GetByEmail) (responses.GetByEmail, error)
}
