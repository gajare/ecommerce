package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	initDB()
}

func initDB() {
	var err error
	connStr := "user=user password=password dbname=ecommerce host=postgres sslmode=disable"
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

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user User
	sqlStatement := `SELECT id, name, email FROM users WHERE id=$1`
	row := db.QueryRow(sqlStatement, id)
	switch err := row.Scan(&user.ID, &user.Name, &user.Email); err {
	case sql.ErrNoRows:
		http.Error(w, "User not found", http.StatusNotFound)
	case nil:
		json.NewEncoder(w).Encode(user)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`SELECT id, name, email FROM users`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}
	json.NewEncoder(w).Encode(users)
}
