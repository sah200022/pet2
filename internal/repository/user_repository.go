package repository

import "errors"

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRepository struct {
	users map[string]User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[string]User),
	}
}

func (r *UserRepository) GetMail(email string) (User, error) {
	user, ok := r.users[email]
	if !ok {
		return User{}, errors.New("user not found")
	}
	return user, nil
}

func (r *UserRepository) Create(user User) (User, error) {
	r.users[user.Email] = user
	return user, nil
}
