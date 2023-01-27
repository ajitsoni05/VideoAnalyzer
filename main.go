package main

import (
	"github/poornachandra7707/myboilerplate/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/getstatus/{key}", handlers.GetStatus).Methods("GET")
	router.HandleFunc("/analyze", handlers.AnalyzeVideo).Methods("POST")
	//handlers.GetStatus()
	log.Fatal(http.ListenAndServe(":4000", router))
}
