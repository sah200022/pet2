package service

import (
	"PetProject/internal/repository"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var jwtSecret = []byte("secret jwt key")

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) Register(email, password string) error {
	if email == "" {
		return errors.New("email is empty")
	}
	if password == "" {
		return errors.New("Password is empty")
	}
	_, err := s.userRepo.GetMail(email)
	if err == nil {
		return errors.New("email already exists")
	}
	user := repository.User{
		Email:    email,
		Password: password,
	}
	_, err = s.userRepo.Create(user)
	return err
}

func (s *AuthService) Login(email, password string) (string, error) {
	if email == "" || password == "" {
		return "", errors.New("email or password is empty")
	}

	user, err := s.userRepo.GetMail(email)
	if err != nil {
		return "", errors.New("email not found")
	}

	if user.Password != password {
		return "", errors.New("incorrect password")
	}

	token, err := s.GenerateToken(email)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) GenerateToken(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Minute * 15).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
