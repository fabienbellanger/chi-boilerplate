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
	Login(requests.UserLogin) (responses.UserLoginRepository, error)
	Create(requests.UserCreationRepository) error
	GetByID(requests.UserByID) (responses.UserByIdRepository, error)
	// GetAll(page, limit, sorts string) (users []entities.User, total int64, err error)
	// GetByEmail(email string) (entities.User, error)
	Delete(requests.UserDelete) error
	// Update(user *entities.User) error
}
