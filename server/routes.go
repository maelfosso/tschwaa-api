package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"tschwaa.com/api/handlers"
	"tschwaa.com/api/services"
)

type signupperMock struct{}

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewStatusResponseWriter(w http.ResponseWriter) *statusResponseWriter {
	return &statusResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (srw *statusResponseWriter) WriteHeader(statusCode int) {
	srw.statusCode = statusCode
	srw.ResponseWriter.WriteHeader(statusCode)
}

func (s *Server) setupRoutes() {
	s.mux.Use(s.requestLoggerMiddleware)

	s.mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // http://localhost:3000", "http://www.tschwaa.local"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	s.mux.Use(services.Verifier)
	s.mux.Use(services.ParseJWTToken)

	s.mux.Group(func(r chi.Router) {
		r.Use(services.Authenticator)
		r.Use(s.convertJwtTokenToMember)

		// Organization
		r.Route("/orgs", func(r chi.Router) {
			handlers.CreateOrganization(r, s.database)
			handlers.ListOrganizations(r, s.database)

			r.Route("/{orgID}", func(r chi.Router) {
				handlers.GetOrganization(r, s.database)
				handlers.GetOrganizationMembers(r, s.database)
				handlers.InviteMembersIntoOrganization(r, s.database)
			})
		})
	})

	s.mux.Group(func(r chi.Router) {
		r.Use(s.convertJwtTokenToMember)

		handlers.Health(s.mux)

		r.Route("/auth/", func(r chi.Router) {
			handlers.Signup(r, s.database)
			handlers.Signin(r, s.database)
		})

		r.Route("/join/", func(r chi.Router) {
			handlers.GetInvitation(r, s.database)
			handlers.JoinOrganization(r, s.database)
		})

	})
}
