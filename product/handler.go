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

func CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	var product Product
	json.NewDecoder(r.Body).Decode(&product)

	sqlStatement := `INSERT INTO products (name, description, price) VALUES ($1, $2, $3) RETURNING id`
	id := 0
	err := db.QueryRow(sqlStatement, product.Name, product.Description, product.Price).Scan(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	product.ID = id
	json.NewEncoder(w).Encode(product)
}

func GetProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var product Product
	sqlStatement := `SELECT id, name, description, price FROM products WHERE id=$1`
	row := db.QueryRow(sqlStatement, id)
	switch err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Price); err {
	case sql.ErrNoRows:
		http.Error(w, "Product not found", http.StatusNotFound)
	case nil:
		json.NewEncoder(w).Encode(product)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`SELECT id, name, description, price FROM products`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	products := []Product{}
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}
	json.NewEncoder(w).Encode(products)
}
