package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"regexp"

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

type UserHandler struct {
	repo *UserRepository
}

var (
	UserRe = regexp.MustCompile(`^/api/users/*$`)
	UserReWithID = regexp.MustCompile(`^/api/users/([a-zA-Z0-9]+(?:-[a-zA-Z0-9]+)+)$`)
)


func (h *UserHandler) Handle (w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	switch {
	case r.Method == http.MethodPost && UserRe.MatchString(path):
		h.Create(w, r)
		return
	case r.Method == http.MethodGet && UserRe.MatchString(path):
		h.GetAll(w, r)
		return
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func throw500(err error, w http.ResponseWriter) {
	w.WriteHeader(500)
	w.Write([]byte(err.Error()))
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.repo.SelectAll()

	if err != nil {
		throw500(err, w)
		return
	}

	res, err := json.Marshal(users)
	if err != nil {
		throw500(err, w)
		return
	}

	w.WriteHeader(200)
	w.Write(res)
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) Insert(username, email, firstName, lastName, avatarUrl, password string) (*User, error) {
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

func (r *UserRepository) SelectByUsername(username string) (*User, error) {
	var user User
	err := r.db.QueryRow("SELECT * FROM users WHERE username = ?", username).Scan(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) SelectAll() ([]User, error) {
	rows, err := r.db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Email, &user.LastName, &user.FirstName, &user.LastName, &user.Email, &user.Password)
		if err != nil{
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil

}



func (h *UserHandler) GetOne(username string) (*User, error) {
	return h.repo.SelectByUsername(username)
}

// func (h *UserHandler) Post(w http.ResponseWriter, r *http.Request) (*User, error) {
// 	r.Body
// 	return h.repo.Insert(username, email, firstName, lastName, avatarUrl, password)
// }

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
