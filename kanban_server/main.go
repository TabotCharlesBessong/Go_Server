package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

type RouteResponse struct {
	Message string `json:message`
	ID      string `json:"id,omitempty"`
}

func main() {
	log.Println("Starting server")

	router := mux.NewRouter()

	log.Println("Setting up routes")

	router.Handle("/register", alice.New(loggingMiddleware).ThenFunc(register)).Methods("POST")
	
	router.Handle("/projects", alice.New(loggingMiddleware).ThenFunc(createPost)).Methods("POST")
	
	router.Handle("/projects/{id}", alice.New(loggingMiddleware).ThenFunc(updatePost)).Methods("PUT")
	
	router.Handle("/projects/{id}", alice.New(loggingMiddleware).ThenFunc(deletePost)).Methods("DELETE")
	
	router.Handle("/login", alice.New(loggingMiddleware).ThenFunc(login)).Methods("POST")
	
	router.Handle("/projects/{id}", alice.New(loggingMiddleware).ThenFunc(getPost)).Methods("GET")
	
	router.Handle("/projects", alice.New(loggingMiddleware).ThenFunc(getPosts)).Methods("GET")
	

	log.Println("Listing on port 5000...")
	log.Fatal(http.ListenAndServe(":5000", router))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

// register
func register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	json.NewEncoder(w).Encode(RouteResponse{Message: "Hello people!"})
}

// register
func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	json.NewEncoder(w).Encode(RouteResponse{Message: "Hello people!"})
}

// login
func createPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	json.NewEncoder(w).Encode(RouteResponse{Message: "Hello people!"})
}

// create post
func updatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	w.Header().Set("Content-type", "application/json")

	json.NewEncoder(w).Encode(RouteResponse{Message: "Hello people!", ID: id})
}

// get post
func getPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	w.Header().Set("Content-type", "application/json")

	json.NewEncoder(w).Encode(RouteResponse{Message: "Hello people!", ID: id})
}

// get posts
func getPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	json.NewEncoder(w).Encode(RouteResponse{Message: "Hello people!"})
}

// delete post
func deletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	w.Header().Set("Content-type", "application/json")

	json.NewEncoder(w).Encode(RouteResponse{Message: "Hello people!", ID: id})
}
