package auth

import (
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

type AuthSessionResponse struct {
	AuthSession struct {
		ID          string  `json:"id" form:"id"`
		RedirectURL string  `json:"redirect_url" form:"redirect_url"`
		AccessID    string  `json:"access_id" form:"access_id"`
		Access      *string `json:"access" form:"access"`
		ExpireAt    string  `json:"expire_at" form:"expire_at"`
		VerifiedAt  *string `json:"verified_at" form:"verified_at"`
		UpdatedAt   string  `json:"updated_at" form:"updated_at"`
		CreatedAt   string  `json:"created_at" form:"created_at"`
	} `json:"auth_session" form:"auth_session" validate:"required,min=8"`
}

type SessionTokenResponse struct {
	AccessToken  string `json:"access_token" form:"access_token"`
	RefreshToken string `json:"refresh_token" form:"refresh_token"`
	TokenType    string `json:"token_type" form:"token_type"`
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
