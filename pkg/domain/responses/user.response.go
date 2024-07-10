package responses

import (
	"chi_boilerplate/pkg/domain/entities"
	vo "chi_boilerplate/pkg/domain/value_objects"
	"time"
)

// UserLogin login response
type UserLogin struct {
	ID        entities.UserID `json:"id" xml:"id"`
	Email     vo.Email        `json:"email" xml:"email"`
	Lastname  string          `json:"lastname" xml:"lastname"`
	Firstname string          `json:"firstname" xml:"firstname"`
	CreatedAt time.Time       `json:"created_at" xml:"created_at"`
	Token     string          `json:"token" xml:"token"`
	ExpiresAt time.Time       `json:"expires_at" xml:"expires_at"`
}

// UserLoginRepository repository login response
type UserLoginRepository struct {
	ID        entities.UserID
	Email     vo.Email
	Lastname  string
	Firstname string
	CreatedAt time.Time
}

// ToUser converts UserLoginRepository to User
func (ulr *UserLoginRepository) ToUser() entities.User {
	return entities.User{
		ID:        ulr.ID,
		Email:     ulr.Email,
		Lastname:  ulr.Lastname,
		Firstname: ulr.Firstname,
	}
}

// UserCreation login response
type UserCreation struct {
	ID        entities.UserID `json:"id" xml:"id"`
	Email     string          `json:"email" xml:"email"`
	Lastname  string          `json:"lastname" xml:"lastname"`
	Firstname string          `json:"firstname" xml:"firstname"`
	CreatedAt time.Time       `json:"created_at" xml:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" xml:"updated_at"`
}

// UsersListPaginated response
type UsersListPaginated struct {
	Data  []entities.User `json:"data"`
	Total int64           `json:"total"`
}
