package entities

import (
	"time"

	vo "chi_boilerplate/pkg/domain/value_objects"
)

// UserID is a type for user ID
type UserID = vo.ID

// User is a struct that represents a user
type User struct {
	ID        UserID      `json:"id" xml:"id" form:"id" validate:"required,uuid"`
	Email     vo.Email    `json:"email" xml:"email" form:"email" validate:"required"`
	Password  vo.Password `json:"-" xml:"-" form:"password" validate:"required"`
	Lastname  string      `json:"lastname" xml:"lastname" form:"lastname" validate:"required"`
	Firstname string      `json:"firstname" xml:"firstname" form:"firstname" validate:"required"`
	CreatedAt time.Time   `json:"created_at" xml:"created_at" form:"created_at"`
	UpdatedAt time.Time   `json:"updated_at" xml:"updated_at" form:"updated_at"`
	DeletedAt *time.Time  `json:"-" xml:"-" form:"deleted_at"`
}
