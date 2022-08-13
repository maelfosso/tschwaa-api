package server

import (
	"tschwaa.com/api/handlers"
)

type signupperMock struct{}

func (s *Server) setupRoutes() {
	handlers.Health(s.mux)

	// Auth
	handlers.Signup(s.mux, s.database)

	// Organization
	handlers.CreateOrganization(s.mux, s.database)
}
