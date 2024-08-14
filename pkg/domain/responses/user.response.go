package responses

import (
	"chi_boilerplate/pkg/domain/entities"
	vo "chi_boilerplate/pkg/domain/value_objects"
	"time"
)

// UsersListPaginated response
type UsersListPaginated struct {
	Data  []entities.User `json:"data"`
	Total int64           `json:"total"`
}

// UserHTTP HTTP response
type UserHTTP struct {
	ID        string `json:"id" xml:"id"`
	Email     string `json:"email" xml:"email"`
	Lastname  string `json:"lastname" xml:"lastname"`
	Firstname string `json:"firstname" xml:"firstname"`
	CreatedAt string `json:"created_at" xml:"created_at"`
	UpdatedAt string `json:"updated_at" xml:"updated_at"`
}

// ======== Login ========

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

// ======== User creation ========

// UserCreation response to create a user
type UserCreation struct {
	ID        entities.UserID `json:"id" xml:"id"`
	Email     vo.Email        `json:"email" xml:"email"`
	Lastname  string          `json:"lastname" xml:"lastname"`
	Firstname string          `json:"firstname" xml:"firstname"`
	CreatedAt time.Time       `json:"created_at" xml:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" xml:"updated_at"`
}

// ToUserHTTP converts UserCreation to UserHTTP
func (u *UserCreation) ToUserHTTP() UserHTTP {
	return UserHTTP{
		ID:        u.ID.String(),
		Email:     u.Email.String(),
		Lastname:  u.Lastname,
		Firstname: u.Firstname,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
	}
}

// ======== Get one user ========

// UserByIDRepository request to get a user by ID
type UserByIdRepository struct {
	ID        string
	Email     string
	Lastname  string
	Firstname string
	CreatedAt string
	UpdatedAt string
}

// UserByID request to get a user by ID
type UserById struct {
	ID        entities.UserID `json:"id" xml:"id"`
	Email     vo.Email        `json:"email" xml:"email"`
	Lastname  string          `json:"lastname" xml:"lastname"`
	Firstname string          `json:"firstname" xml:"firstname"`
	CreatedAt time.Time       `json:"created_at" xml:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" xml:"updated_at"`
}

// ToUserHTTP converts UserById to UserHTTP
func (u *UserById) ToUserHTTP() UserHTTP {
	return UserHTTP{
		ID:        u.ID.String(),
		Email:     u.Email.String(),
		Lastname:  u.Lastname,
		Firstname: u.Firstname,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
	}
}

// ======== Get all users ========

type UsersList Pagination[UsersListModel]

// UsersListModel request a user by ID
type UsersListModel struct {
	ID        string `json:"id" xml:"id"`
	Email     string `json:"email" xml:"email"`
	Lastname  string `json:"lastname" xml:"lastname"`
	Firstname string `json:"firstname" xml:"firstname"`
	CreatedAt string `json:"created_at" xml:"created_at"`
	UpdatedAt string `json:"updated_at" xml:"updated_at"`
}
