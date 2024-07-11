package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
)

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/health", HealthCheckHandler).Methods("GET")
    r.HandleFunc("/products", CreateProductHandler).Methods("POST")
    r.HandleFunc("/products/{id}", GetProductHandler).Methods("GET")
    r.HandleFunc("/products", GetProductsHandler).Methods("GET")

    log.Println("Starting Product service on :8002")
    if err := http.ListenAndServe(":8002", r); err != nil {
        log.Fatalf("could not start server: %s\n", err)
    }
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}
