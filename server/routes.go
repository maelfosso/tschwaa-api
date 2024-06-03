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
		AllowedOrigins:   []string{"http://localhost:5173", "http://www.tschwaa.local"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	s.mux.Use(services.Verifier)
	s.mux.Use(services.ParseJWTToken)

	// Protected Routes
	s.mux.Group(func(r chi.Router) {
		r.Use(services.Authenticator)
		r.Use(s.convertJWTTokenToMember)

		handlers.GetCurrentUser(r)

		// Organization
		r.Route("/orgs", func(r chi.Router) {
			handlers.CreateOrganization(r, s.database.Storage)
			handlers.ListOrganizations(r, s.database.Storage)

			r.Route("/{orgID}", func(r chi.Router) {

				r.Route("/sessions", func(r chi.Router) {
					handlers.CreateSession(r, s.database.Storage)
					handlers.GetCurrentSession(r, s.database.Storage)

					r.Route("/{sessionID}", func(r chi.Router) {

						r.Route("/members", func(r chi.Router) {
							handlers.GetMembersOfSession(r, s.database.Storage)
							handlers.AddMemberToSession(r, s.database.Storage)
							handlers.UpdateSessionMembers(r, s.database.Storage)
							handlers.RemoveMemberFromSession(r, s.database.Storage)
						})

						r.Route("/places", func(r chi.Router) {
							handlers.GetPlaceOfSession(r, s.database.Storage)
							handlers.UpdatePlaceOfSession(r, s.database.Storage)
							handlers.ChangePlaceOfSession(r, s.database.Storage)
						})

					})
				})

				handlers.GetOrganization(r, s.database.Storage)
				handlers.GetOrganizationMembers(r, s.database.Storage)
				handlers.InviteMembersIntoOrganization(r, s.database.Storage)
			})
		})
	})

	// Public Route
	s.mux.Group(func(r chi.Router) {
		r.Use(s.convertJWTTokenToMember)

		handlers.Health(s.mux)

		r.Route("/auth/", func(r chi.Router) {
			handlers.SignUp(r, s.database.Storage)
			handlers.SignIn(r, s.database.Storage)
			handlers.GetOtp(r, s.database.Storage)
			handlers.CheckOtp(r, s.database.Storage)
			handlers.ResendOtp(r, s.database.Storage)
		})

		r.Route("/token", func(r chi.Router) {
			handlers.Refresh(r)
			handlers.IsTokenValid(r)
		})

		r.Route("/join/", func(r chi.Router) {
			handlers.GetInvitation(r, s.database.Storage)
			handlers.JoinOrganization(r, s.database.Storage)
		})

	})
}
