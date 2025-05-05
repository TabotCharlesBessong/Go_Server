package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type RouteResponse struct {
	Message string `json:message`
}

func main() {
	fmt.Println("Hello world")

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter,r *http.Request) {
		w.Header().Set("Content-type","application/json")

		json.NewEncoder(w).Encode(RouteResponse{Message: "Hello people!"})
	}).Methods("GET")

	log.Fatal(http.ListenAndServe(":5000",router))
}