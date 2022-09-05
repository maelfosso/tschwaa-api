package server

import (
	"github.com/go-chi/cors"
	"tschwaa.com/api/handlers"
)

type signupperMock struct{}

func (s *Server) setupRoutes() {
	s.mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	handlers.Health(s.mux)

	// Auth
	handlers.Signup(s.mux, s.database)
	handlers.Signin(s.mux, s.database)

	// Organization
	handlers.CreateOrganization(s.mux, s.database)
}
