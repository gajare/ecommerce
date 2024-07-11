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
    var err error
    connStr := "user=user password=password dbname=ecommerce sslmode=disable"
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
}

func CreateCartItemHandler(w http.ResponseWriter, r *http.Request) {
    var item CartItem
    if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    sqlStatement := `INSERT INTO cart (user_id, product_id, quantity) VALUES ($1, $2, $3) RETURNING id`
    var id int
    if err := db.QueryRow(sqlStatement, item.UserID, item.ProductID, item.Quantity).Scan(&id); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    item.ID = id
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(item)
}

func GetCartItemHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    var item CartItem
    sqlStatement := `SELECT id, user_id, product_id, quantity FROM cart WHERE id = $1`
    row := db.QueryRow(sqlStatement, id)
    switch err := row.Scan(&item.ID, &item.UserID, &item.ProductID, &item.Quantity); err {
    case sql.ErrNoRows:
        http.Error(w, "Cart item not found", http.StatusNotFound)
        return
    case nil:
        json.NewEncoder(w).Encode(item)
    default:
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func UpdateCartItemHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    var item CartItem
    if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    sqlStatement := `UPDATE cart SET user_id = $1, product_id = $2, quantity = $3 WHERE id = $4`
    res, err := db.Exec(sqlStatement, item.UserID, item.ProductID, item.Quantity, id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    rowsAffected, err := res.RowsAffected()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if rowsAffected == 0 {
        http.Error(w, "Cart item not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func DeleteCartItemHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    sqlStatement := `DELETE FROM cart WHERE id = $1`
    res, err := db.Exec(sqlStatement, id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    rowsAffected, err := res.RowsAffected()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if rowsAffected == 0 {
        http.Error(w, "Cart item not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
