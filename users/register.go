package users

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type UserInput struct {
	Email    string `json:"email"`
	UserName string `json:"userName"`
	Password string `json:"password"`
}

// encore:api public method=POST path=/register
func (s Service) Register(context context.Context, user *UserInput) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	site := &User{
		Email:          user.Email,
		UserName:       user.UserName,
		HashedPassword: string(hashed),
	}
	if err := s.db.Create(site).Error; err != nil {
		return err
	}
	return nil
}
