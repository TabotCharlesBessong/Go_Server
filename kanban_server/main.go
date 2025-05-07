package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"kanban_server/config"
	"kanban_server/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

type RouteResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type ProjectRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

var jwtKey = []byte("your-secret-key") // In production, use environment variable

func main() {
	log.Println("Starting server")

	// Initialize database
	config.InitDB()

	router := mux.NewRouter()

	log.Println("Setting up routes")

	router.Handle("/register", alice.New(loggingMiddleware).ThenFunc(register)).Methods("POST")
	router.Handle("/login", alice.New(loggingMiddleware).ThenFunc(login)).Methods("POST")
	router.Handle("/projects", alice.New(loggingMiddleware, authMiddleware).ThenFunc(createProject)).Methods("POST")
	router.Handle("/projects/{id}", alice.New(loggingMiddleware, authMiddleware).ThenFunc(updateProject)).Methods("PUT")
	router.Handle("/projects/{id}", alice.New(loggingMiddleware, authMiddleware).ThenFunc(deleteProject)).Methods("DELETE")
	router.Handle("/projects/{id}", alice.New(loggingMiddleware, authMiddleware).ThenFunc(getProject)).Methods("GET")
	router.Handle("/projects", alice.New(loggingMiddleware, authMiddleware).ThenFunc(getProjects)).Methods("GET")

	log.Println("Listening on port 5000...")
	log.Fatal(http.ListenAndServe(":5000", router))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			json.NewEncoder(w).Encode(RouteResponse{Error: "Authorization header required"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			json.NewEncoder(w).Encode(RouteResponse{Error: "Invalid token"})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(RouteResponse{Error: "Invalid request body"})
		return
	}

	user := models.User{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		json.NewEncoder(w).Encode(RouteResponse{Error: "Failed to create user"})
		return
	}

	json.NewEncoder(w).Encode(RouteResponse{
		Message: "User registered successfully",
		Data:    user,
	})
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(RouteResponse{Error: "Invalid request body"})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		json.NewEncoder(w).Encode(RouteResponse{Error: "Invalid credentials"})
		return
	}

	if err := user.CheckPassword(req.Password); err != nil {
		json.NewEncoder(w).Encode(RouteResponse{Error: "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		json.NewEncoder(w).Encode(RouteResponse{Error: "Failed to generate token"})
		return
	}

	json.NewEncoder(w).Encode(RouteResponse{
		Message: "Login successful",
		Data:    map[string]string{"token": tokenString},
	})
}

func createProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req ProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(RouteResponse{Error: "Invalid request body"})
		return
	}

	project := models.Project{
		Title:       req.Title,
		Description: req.Description,
		Status:      "active",
	}

	if err := config.DB.Create(&project).Error; err != nil {
		json.NewEncoder(w).Encode(RouteResponse{Error: "Failed to create project"})
		return
	}

	json.NewEncoder(w).Encode(RouteResponse{
		Message: "Project created successfully",
		Data:    project,
	})
}

func updateProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]

	var req ProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(RouteResponse{Error: "Invalid request body"})
		return
	}

	var project models.Project
	if err := config.DB.First(&project, id).Error; err != nil {
		json.NewEncoder(w).Encode(RouteResponse{Error: "Project not found"})
		return
	}

	project.Title = req.Title
	project.Description = req.Description

	if err := config.DB.Save(&project).Error; err != nil {
		json.NewEncoder(w).Encode(RouteResponse{Error: "Failed to update project"})
		return
	}

	json.NewEncoder(w).Encode(RouteResponse{
		Message: "Project updated successfully",
		Data:    project,
	})
}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]

	if err := config.DB.Delete(&models.Project{}, id).Error; err != nil {
		json.NewEncoder(w).Encode(RouteResponse{Error: "Failed to delete project"})
		return
	}

	json.NewEncoder(w).Encode(RouteResponse{
		Message: "Project deleted successfully",
	})
}

func getProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]

	var project models.Project
	if err := config.DB.Preload("Tasks").First(&project, id).Error; err != nil {
		json.NewEncoder(w).Encode(RouteResponse{Error: "Project not found"})
		return
	}

	json.NewEncoder(w).Encode(RouteResponse{
		Data: project,
	})
}

func getProjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var projects []models.Project
	if err := config.DB.Preload("Tasks").Find(&projects).Error; err != nil {
		json.NewEncoder(w).Encode(RouteResponse{Error: "Failed to fetch projects"})
		return
	}

	json.NewEncoder(w).Encode(RouteResponse{
		Data: projects,
	})
}
