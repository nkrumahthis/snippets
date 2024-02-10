package main

import "database/sql"

type Database struct {
	db *sql.DB
}

func (db *Database) Init() error{
	_, err := db.db.Exec(`
		PRAGMA foreign_keys = ON;
		CREATE TABLE IF NOT EXISTS users (
			id VARCHAR(30) PRIMARY KEY,
			username VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL,
			first_name VARCHAR(255) NOT NULL,
			last_name VARCHAR(255) NOT NULL,
			avatar_url VARCHAR(512) NOT NULL,
			password VARCHAR(255) NOT NULL
		);

		CREATE TABLE IF NOT EXISTS snippets(
			id VARCHAR(30) PRIMARY KEY,
			user_id VARCHAR(30) NOT NULL FOREIGN KEY (id) REFERENCES users (id) ON DELETE CASCADE,
			title TEXT,
			code TEXT,
			upvotes INTEGER DEFAULT 0,
			downvotes INTEGER DEFAULT 0,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS upvotes (
			id VARCHAR(30) PRIMARY KEY,
			user_id VARCHAR(30) NOT NULL FOREIGN KEY (id) REFERENCES users (id) ON DELETE CASCADE,
			timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE TABLE IF NOT EXISTS downvotes(
			id VARCHAR(30) PRIMARY KEY,
			user_id VARCHAR(30) NOT NULL FOREIGN KEY (id) REFERENCES users (id) ON DELETE CASCADE,
			timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`)
	return err
}
