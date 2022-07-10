package server

import "tschwaa.com/api/handlers"

func (s *Server) setupRoutes() {
	handlers.Health(s.mux)
}
