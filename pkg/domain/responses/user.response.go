package responses

import (
	"chi_boilerplate/pkg/domain/entities"
	vo "chi_boilerplate/pkg/domain/value_objects"
	"time"
)

// UserLogin login response
type UserLogin struct {
	ID        entities.UserID `json:"id" xml:"id"`
	Email     string          `json:"email" xml:"email"`
	Lastname  string          `json:"lastname" xml:"lastname"`
	Firstname string          `json:"firstname" xml:"firstname"`
	CreatedAt string          `json:"created_at" xml:"created_at"`
	Token     string          `json:"token" xml:"token"`
	ExpiresAt string          `json:"expires_at" xml:"expires_at"`
}

// UserLoginRepository repository login response
type UserLoginRepository struct {
	ID        entities.UserID
	Email     string
	Lastname  string
	Firstname string
	CreatedAt string
}

// ToUser converts UserLoginRepository to User
// TODO: Add tests
func (ulr *UserLoginRepository) ToUser() (entities.User, error) {
	email, err := vo.NewEmail(ulr.Email)
	if err != nil {
		return entities.User{}, err
	}

	createdAt, err := time.Parse(time.RFC3339, ulr.CreatedAt)
	if err != nil {
		return entities.User{}, err
	}

	return entities.User{
		ID:        ulr.ID,
		Email:     email,
		Lastname:  ulr.Lastname,
		Firstname: ulr.Firstname,
		CreatedAt: createdAt,
	}, nil
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
