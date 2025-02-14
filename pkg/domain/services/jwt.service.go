package services

import (
	"chi_boilerplate/pkg/domain/entities"
	"chi_boilerplate/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

// JWT reprensents a JWT token
type JWT struct {
	Value     string
	ExpiredAt time.Time
}

// NewJWT creates a new JWT token
func NewJWT(id entities.UserID, lifetime time.Duration, algo, secret string) (JWT, error) {
	// Create token and key
	token, key, err := utils.GetTokenAndKeyFromAlgo(algo, secret, viper.GetString("JWT_PRIVATE_KEY_PATH"))
	if err != nil {
		return JWT{}, err
	}

	// Expiration time
	now := time.Now()
	expiresAt := now.Add(time.Hour * lifetime)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = id.String()
	claims["exp"] = expiresAt.Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	// Generate encoded token and send it as response
	t, err := token.SignedString(key)
	if err != nil {
		return JWT{}, err
	}
	return JWT{Value: t, ExpiredAt: expiresAt}, nil
}
