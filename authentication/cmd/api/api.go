package api

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
}

func NewAPIServer(addr string) *APIServer {
	return &APIServer{
		addr: addr,
	}
}

func (s *APIServer) Run() error {
	// create a router
	router := mux.NewRouter()

	// register services

	c := cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PATCH, DELETE",
		AllowCredentials: true,
	})
	log.Println("Listing on port",s.addr)

	return http.ListenAndServe(s.addr, router)
}