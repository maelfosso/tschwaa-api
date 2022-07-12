package server

import (
	"context"

	"tschwaa.com/api/handlers"
	"tschwaa.com/api/model"
)

type signupperMock struct{}

func (s *Server) setupRoutes() {
	handlers.Health(s.mux)

	handlers.Signup(s.mux, &signupperMock{})
}

func (s signupperMock) Signup(ctx context.Context, user model.User) error {
	return nil
}
