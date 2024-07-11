package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/health", HealthCheckHandler).Methods("GET")
	r.HandleFunc("/cart", CreateCartItemHandler).Methods("POST")
	r.HandleFunc("/cart/{id}", GetCartItemHandler).Methods("GET")
	r.HandleFunc("/cart/{id}", UpdateCartItemHandler).Methods("PUT")
	r.HandleFunc("/cart/{id}", DeleteCartItemHandler).Methods("DELETE")

	log.Println("Starting Cart service on :8004")
	if err := http.ListenAndServe(":8004", r); err != nil {
		log.Fatalf("could not start server: %s\n", err)
	}
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
