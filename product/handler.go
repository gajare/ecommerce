package main

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
    "github.com/gorilla/mux"
    _ "github.com/lib/pq"
    "strconv"
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
    sqlStatement := `SELECT id, name, description, price FROM products WHERE id = $1`
    row := db.QueryRow(sqlStatement, id)
    switch err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Price); err {
    case sql.ErrNoRows:
        http.Error(w, "Product not found", http.StatusNotFound)
        return
    case nil:
        json.NewEncoder(w).Encode(product)
    default:
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func GetAllProductsHandler(w http.ResponseWriter, r *http.Request) {
    var products []Product

    rows, err := db.Query("SELECT id, name, description, price FROM products")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    for rows.Next() {
        var product Product
        if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        products = append(products, product)
    }

    if err := rows.Err(); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(products)
}

func UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid product ID", http.StatusBadRequest)
        return
    }

    var product Product
    json.NewDecoder(r.Body).Decode(&product)

    sqlStatement := `UPDATE products SET name = $1, description = $2, price = $3 WHERE id = $4`
    res, err := db.Exec(sqlStatement, product.Name, product.Description, product.Price, id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    count, err := res.RowsAffected()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if count == 0 {
        http.Error(w, "Product not found", http.StatusNotFound)
        return
    }

    product.ID = id
    json.NewEncoder(w).Encode(product)
}

func DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid product ID", http.StatusBadRequest)
        return
    }

    sqlStatement := `DELETE FROM products WHERE id = $1`
    res, err := db.Exec(sqlStatement, id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    count, err := res.RowsAffected()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if count == 0 {
        http.Error(w, "Product not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
