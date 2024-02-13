package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func init() {

}

func main() {
	db, err := sql.Open("sqlite3", "temp/data.db")
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	defer db.Close()

	database := Database{db}
	if err := database.Init(); err != nil {
		log.Fatal("Error initializing database:", err)
	}

	userHandler := UserHandler{NewUserRepository(db)}
	
	fileServer := http.FileServer(http.Dir("./_ui/dist"))
	mux := http.NewServeMux()
	mux.Handle("/", fileServer)
	mux.HandleFunc("/api/users", userHandler.Handle)

	print("Starting server on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Println("server failed:", err)
	}
}
