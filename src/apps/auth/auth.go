package auth

import (
	"math"
	"math/rand/v2"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type RegisterForm struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Username  *string `json:"username"`
	Email     string  `json:"email" validate:"required,email"`
	Password  *string `json:"password"`
}

type LoginForm struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=8"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func GenerateFullTokens(id string) (map[string]any, error) {
	accessToken, err := GenerateToken(id, false)
	if err != nil {
		return nil, err
	}
	refreshToken, err := GenerateToken(id, true)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
	}, nil
}

func GenerateUsername(email string) string {
	var username string = email
	var re *regexp.Regexp

	re = regexp.MustCompile("@.*$")
	username = re.ReplaceAllString(username, "")

	re = regexp.MustCompile("[^a-z0-9._-]")
	username = re.ReplaceAllString(username, "-")

	re = regexp.MustCompile("[._-]{2,}")
	username = re.ReplaceAllString(username, "-")

	username = strings.ToLower(username)
	username = username[0:int(math.Min(float64(len(username)), 20))]

	username = username + strconv.Itoa(int(1000+rand.Float64()*9000))

	return username
}
