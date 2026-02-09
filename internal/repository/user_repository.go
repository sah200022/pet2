package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID       int    `db:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetMail(email string) (User, error) {
	query := `
	SELECT id, email, password 
	FROM users
	WHERE email = $1;
`
	var user User

	err := r.db.QueryRow(context.Background(), query, email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (r *UserRepository) Create(user User) (User, error) {
	query := `
		INSERT INTO users (email, password)
		VALUES ($1, $2)
		RETURNING id
	`
	var id int

	err := r.db.QueryRow(context.Background(), query, user.Email, user.Password).Scan(&id)
	if err != nil {
		return User{}, err
	}
	user.ID = id
	return user, nil
}
