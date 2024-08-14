package requests

// UserLogin request
type UserLogin struct {
	Email    string `json:"email" xml:"email" form:"email" validate:"required,email"`
	Password string `json:"password" xml:"password" form:"password" validate:"required,min=8"`
}

// UserByID request
type UserByID struct {
	ID string `json:"id" xml:"id" form:"id" validate:"required,uuid"`
}

// UserCreation request to create a user
type UserCreation struct {
	Email     string `json:"email" xml:"email" form:"email" validate:"required,email"`
	Password  string `json:"password" xml:"password" form:"password" validate:"required,min=8"`
	Lastname  string `json:"lastname" xml:"lastname" form:"lastname" validate:"required"`
	Firstname string `json:"firstname" xml:"firstname" form:"firstname" validate:"required"`
}

// UserUpdate request to update a user
type UserUpdate struct {
	ID        string `json:"id" xml:"id" form:"id" validate:"required,uuid"`
	Email     string `json:"email" xml:"email" form:"email" validate:"required,email"`
	Password  string `json:"password" xml:"password" form:"password" validate:"required,min=8"`
	Lastname  string `json:"lastname" xml:"lastname" form:"lastname" validate:"required"`
	Firstname string `json:"firstname" xml:"firstname" form:"firstname" validate:"required"`
}

// UserCreation request to create a user
type UserCreationRepository struct {
	ID        string
	Email     string
	Password  string
	Lastname  string
	Firstname string
	CreatedAt string
	UpdatedAt string
}

// UserDelete request
type UserDelete struct {
	ID string `json:"id" xml:"id" form:"id" validate:"required,uuid"`
}

// UsersList request
type UsersList Pagination
