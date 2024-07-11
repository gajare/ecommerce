package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
)

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/health", HealthCheckHandler).Methods("GET")
    r.HandleFunc("/users", CreateUserHandler).Methods("POST")
    r.HandleFunc("/users/{id}", GetUserHandler).Methods("GET")
    r.HandleFunc("/users", GetUsersHandler).Methods("GET")

    log.Println("Starting User service on :8001")
    if err := http.ListenAndServe(":8001", r); err != nil {
        log.Fatalf("could not start server: %s\n", err)
    }
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}
