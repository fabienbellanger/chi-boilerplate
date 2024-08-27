package entities

import (
	"chi_boilerplate/utils"
	"crypto/sha512"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"

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

// HashUserPassword hashes a password
func HashUserPassword(password string) string {
	passwordBytes := sha512.Sum512([]byte(password))
	return hex.EncodeToString(passwordBytes[:])
}

// GenerateJWT returns a token
func (u *User) GenerateJWT(lifetime time.Duration, algo, secret string) (string, time.Time, error) {
	// Create token and key
	token, key, err := utils.GetTokenAndKeyFromAlgo(algo, secret, viper.GetString("JWT_PRIVATE_KEY_PATH"))
	if err != nil {
		return "", time.Now(), err
	}

	// Expiration time
	now := time.Now()
	expiresAt := now.Add(time.Hour * lifetime)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = u.ID.String()
	claims["exp"] = expiresAt.Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	// Generate encoded token and send it as response
	t, err := token.SignedString(key)
	if err != nil {
		return "", expiresAt, err
	}

	return t, expiresAt, nil
}
