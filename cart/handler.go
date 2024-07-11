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

func AddToCartHandler(w http.ResponseWriter, r *http.Request) {
    var cartItem CartItem
    json.NewDecoder(r.Body).Decode(&cartItem)

    sqlStatement := `INSERT INTO cart (user_id, product_id, quantity) VALUES ($1, $2, $3) RETURNING id`
    id := 0
    err := db.QueryRow(sqlStatement, cartItem.UserID, cartItem.ProductID, cartItem.Quantity).Scan(&id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    cartItem.ID = id
    json.NewEncoder(w).Encode(cartItem)
}

func GetCartHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID, err := strconv.Atoi(vars["user_id"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    rows, err := db.Query(`SELECT id, user_id, product_id, quantity FROM cart WHERE user_id=$1`, userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    cartItems := []CartItem{}
    for rows.Next() {
        var cartItem CartItem
        if err := rows.Scan(&cartItem.ID, &cartItem.UserID, &cartItem.ProductID, &cartItem.Quantity); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        cartItems = append(cartItems, cartItem)
    }
    json.NewEncoder(w).Encode(cartItems)
}
