package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	var err error
	connStr := "user=user password=password dbname=ecommerce sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)

	sqlStatement := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`
	id := 0
	err := db.QueryRow(sqlStatement, user.Name, user.Email).Scan(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.ID = id
	json.NewEncoder(w).Encode(user)
}

// Add other handlers here (e.g., GetUserHandler, UpdateUserHandler, etc.)
