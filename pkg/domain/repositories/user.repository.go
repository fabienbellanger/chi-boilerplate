package repositories

import (
	"chi_boilerplate/pkg/domain/entities"
	"errors"
)

var (
	// ErrUserNotFound is the error returned when a user is not found.
	ErrUserNotFound = errors.New("user not found")
)

// UserRepository is the interface that wraps the basic user repository methods.
type UserRepository interface {
	Login(email, password string) (entities.User, error)
	Create(user *entities.User) error
	// GetAll(page, limit, sorts string) (users []entities.User, total int64, err error)
	// GetByID(id string) (entities.User, error)
	// GetByEmail(email string) (entities.User, error)
	// Delete(id string) error
	// Update(user *entities.User) error
}
