package main

import (
	"database/sql"
	"time"

	"github.com/oklog/ulid/v2"
)

type Snippet struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	Title      string    `json:"title"`
	Code       string    `json:"code"`
	Upvotes    int       `json:"upvotes"`
	Downvotes  int       `json:"downvotes"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}

type Downvote struct {
	UserID     string    `json:"user_id"`
	Snipped_id string    `json:"snipped_id"`
	Timestamp  time.Time `json:"timestamp"`
}

type Upvote struct {
	UserID     string    `json:"user_id"`
	Snippet_id string    `json:"snippet_id"`
	Timestamp  time.Time `json:"timestamp"`
}

type SnippetRepository struct {
	db *sql.DB
}

func (r *SnippetRepository) Create(title, code, userId string) (*Snippet, error) {
	var snippet Snippet
	snippet.ID = ulid.Make().String()
	snippet.Title = title
	snippet.Code = code
	snippet.UserID = userId
	
	_, err := r.db.Exec("INSERT INTO snippet (id, title, code, user_id) VALUES (?, ?, ?, ?)", snippet.ID, snippet.Title, snippet.Code, snippet.UserID)
	if err != nil {
		return nil, err
	}
	return &snippet, nil
}

