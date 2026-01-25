package service

import "errors"

type AuthService struct {
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Register(email, password string) error {
	if email == "" {
		return errors.New("email is empty")
	}
	if password == "" {
		return errors.New("Password is empty")
	}
	return nil
}
