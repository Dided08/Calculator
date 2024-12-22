package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/Dided08/Calculator/internal/handlers"
)

func main() {
    r := mux.NewRouter()
    api := r.PathPrefix("/api/v1").Subrouter()
    api.HandleFunc("/calculate", handlers.CalculateHandler).Methods("POST")

    log.Println("Starting server on :8080...")
    log.Fatal(http.ListenAndServe(":8080", r))
}
