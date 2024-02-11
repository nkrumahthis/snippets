package main

import (
	"database/sql"

	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	AvatarUrl string `json:"avatar_url"`
	Password  string `json:"password"`
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) Create(username, email, firstName, lastName, avatarUrl, password string) (*User, error) {

	var user User
	user.ID = ulid.Make().String()
	user.Username = username
	user.Email = email
	user.FirstName = firstName
	user.LastName = lastName
	user.AvatarUrl = avatarUrl
	user.Password, _ = hashPassword(password)

	_, err := r.db.Exec(`INSERT INTO users (id, username, email, first_name, last_name, avatar_url, password) VALUES (?, ?, ?, ?, ?, ?, ?)`, user.ID, user.Username, user.Email, user.FirstName, user.LastName, user.AvatarUrl, user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Get(username string) (*User, error) {
	var user User
	err := r.db.QueryRow("SELECT * FROM users WHERE username = ?", username).Scan(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
