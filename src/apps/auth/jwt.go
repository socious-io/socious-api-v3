package auth

import (
	"errors"
	"socious/src/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	ID      string `json:"id"`
	Refresh bool   `json:"refresh"`
	jwt.RegisteredClaims
}

type SSOClaims struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	jwt.RegisteredClaims
}

func GenerateToken(id string, refresh bool) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		ID:      id,
		Refresh: refresh,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Config.Secret))
}

func GenerateSSOToken(email, firstName, lastName string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &SSOClaims{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Config.SSO.Secret))
}

func VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.Secret), nil
	})
	if err != nil {
		return nil, err
	} else if !token.Valid {
		return nil, errors.New("invalid token")
	} else if claims, ok := token.Claims.(*Claims); ok {
		return claims, nil
	}

	return nil, errors.New("unknown claims type, cannot proceed")
}
