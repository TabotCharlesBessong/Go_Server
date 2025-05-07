package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"kanban_server/config"
	"kanban_server/models"

	"github.com/golang-jwt/jwt/v5"
)

func TestRegister(t *testing.T) {
	// Initialize test database
	config.InitDB()

	// Create test request
	reqBody := RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	register(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check response body
	var response RouteResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if response.Error != "" {
		t.Errorf("handler returned error: %v", response.Error)
	}

	// Verify user was created in database
	var user models.User
	if err := config.DB.Where("email = ?", reqBody.Email).First(&user).Error; err != nil {
		t.Errorf("user was not created in database: %v", err)
	}
}

func TestLogin(t *testing.T) {
	// Initialize test database
	config.InitDB()

	// Create test user
	user := models.User{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}
	config.DB.Create(&user)

	// Create test request
	reqBody := LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	login(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check response body
	var response RouteResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if response.Error != "" {
		t.Errorf("handler returned error: %v", response.Error)
	}

	// Verify token was returned
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Error("response data is not a map")
	}
	if _, ok := data["token"]; !ok {
		t.Error("token not found in response")
	}
}

func TestCreateProject(t *testing.T) {
	// Initialize test database
	config.InitDB()

	// Create test user and get token
	user := models.User{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}
	config.DB.Create(&user)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, _ := token.SignedString(jwtKey)

	// Create test request
	reqBody := ProjectRequest{
		Title:       "Test Project",
		Description: "Test Description",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/projects", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokenString)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	createProject(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check response body
	var response RouteResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if response.Error != "" {
		t.Errorf("handler returned error: %v", response.Error)
	}

	// Verify project was created in database
	var project models.Project
	if err := config.DB.Where("title = ?", reqBody.Title).First(&project).Error; err != nil {
		t.Errorf("project was not created in database: %v", err)
	}
}

func TestGetProjects(t *testing.T) {
	// Initialize test database
	config.InitDB()

	// Create test user and get token
	user := models.User{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}
	config.DB.Create(&user)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, _ := token.SignedString(jwtKey)

	// Create test projects
	projects := []models.Project{
		{Title: "Project 1", Description: "Description 1", Status: "active"},
		{Title: "Project 2", Description: "Description 2", Status: "active"},
	}
	for _, p := range projects {
		config.DB.Create(&p)
	}

	// Create test request
	req := httptest.NewRequest("GET", "/projects", nil)
	req.Header.Set("Authorization", tokenString)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	getProjects(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check response body
	var response RouteResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if response.Error != "" {
		t.Errorf("handler returned error: %v", response.Error)
	}

	// Verify projects were returned
	projectsData, ok := response.Data.([]interface{})
	if !ok {
		t.Error("response data is not a slice")
	}
	if len(projectsData) != 2 {
		t.Errorf("expected 2 projects, got %d", len(projectsData))
	}
}
