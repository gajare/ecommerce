package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
)

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/health", HealthCheckHandler).Methods("GET")
    r.HandleFunc("/purchases", CreatePurchaseHandler).Methods("POST")
    r.HandleFunc("/purchases/{id}", GetPurchaseHandler).Methods("GET")
    r.HandleFunc("/purchases", GetPurchasesHandler).Methods("GET")

    log.Println("Starting Purchase service on :8003")
    if err := http.ListenAndServe(":8003", r); err != nil {
        log.Fatalf("could not start server: %s\n", err)
    }
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}
