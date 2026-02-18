package service

import (
	"PetProject/internal/user/repository"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	JWTSecret []byte
}

func NewAuthService(userRepo *repository.UserRepository, JWT_SECRET []byte) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		JWTSecret: JWT_SECRET,
	}
}

func (s *AuthService) Register(email, password string) error {
	if email == "" {
		return errors.New("email is empty")
	}
	if password == "" {
		return errors.New("Password is empty")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}
	_, err = s.userRepo.GetMail(email)
	if err == nil {
		return errors.New("email already exists")
	}
	user := repository.User{
		Email:    email,
		Password: string(hashedPassword),
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

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)
	if err != nil {
		return "", errors.New("incorrect password")
	}

	token, err := s.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) GenerateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Minute * 15).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.JWTSecret)
}
