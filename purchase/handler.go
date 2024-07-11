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
    connStr := "user=user password=password dbname=ecommerce sslmode=disable"
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
}

func CreatePurchaseHandler(w http.ResponseWriter, r *http.Request) {
    var purchase Purchase
    json.NewDecoder(r.Body).Decode(&purchase)

    sqlStatement := `INSERT INTO purchases (user_id, product_id, quantity, total_price) VALUES ($1, $2, $3, $4) RETURNING id`
    id := 0
    err := db.QueryRow(sqlStatement, purchase.UserID, purchase.ProductID, purchase.Quantity, purchase.TotalPrice).Scan(&id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    purchase.ID = id
    json.NewEncoder(w).Encode(purchase)
}

func GetPurchaseHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid purchase ID", http.StatusBadRequest)
        return
    }

    var purchase Purchase
    sqlStatement := `SELECT id, user_id, product_id, quantity, total_price FROM purchases WHERE id=$1`
    row := db.QueryRow(sqlStatement, id)
    switch err := row.Scan(&purchase.ID, &purchase.UserID, &purchase.ProductID, &purchase.Quantity, &purchase.TotalPrice); err {
    case sql.ErrNoRows:
        http.Error(w, "Purchase not found", http.StatusNotFound)
    case nil:
        json.NewEncoder(w).Encode(purchase)
    default:
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func GetPurchasesHandler(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query(`SELECT id, user_id, product_id, quantity, total_price FROM purchases`)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    purchases := []Purchase{}
    for rows.Next() {
        var purchase Purchase
        if err := rows.Scan(&purchase.ID, &purchase.UserID, &purchase.ProductID, &purchase.Quantity, &purchase.TotalPrice); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        purchases = append(purchases, purchase)
    }
    json.NewEncoder(w).Encode(purchases)
}
