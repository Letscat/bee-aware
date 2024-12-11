package users

import (
	"context"
	"net/http"
	"time"

	"encore.dev/beta/errs"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserLogin struct {
	Email    string
	Password string
}
type LoginResponse struct {
	SessionID string `header:"Set-Cookie"`
}

// encore:api public method=POST path=/login
func (s Service) Login(context context.Context, login *UserLogin) (*LoginResponse, error) {
	var user User
	if err := s.db.Where("email = $1", login.Email).First(&user).Error; err != nil {
		return nil, err
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(login.Password))
	if err != nil {
		return nil, &errs.Error{
			Code:    errs.Unauthenticated,
			Message: err.Error(),
		}
	}
	token, err := createJWTToken("1234")
	if err != nil {
		return nil, &errs.Error{
			Code:    errs.Internal,
			Message: err.Error(),
		}
	}
	return &LoginResponse{SessionID: token}, nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func createJWTToken(uid string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 72).Unix(),
		"iat": time.Now().Unix(),
		"uid": uid,
	})
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		return "", err
	}
	cookie := &http.Cookie{
		Name:     "session",
		Value:    tokenString,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(time.Hour * 72),
	}
	return cookie.String(), nil
}
