package server

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/maelfosso/jwtauth"
	"go.uber.org/zap"
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

	s.mux.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(services.TokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				_, claims, _ := jwtauth.FromContext(req.Context())
				s.log.Info(
					"JWT Claims",
					// zap.String("email", fmt.Sprintf("%v", claims["email"])),
					zap.Any("Jwt claims", claims),
				)

				user, err := s.database.FindUserByEmail(req.Context(), claims["Email"].(string))
				if err != nil {
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}

				s.log.Info(
					"JWT User",
					zap.Any("Jwt user", user),
				)

				ctx := req.Context()
				ctx = context.WithValue(ctx, services.JwtUserKey, user)
				next.ServeHTTP(w, req.WithContext(ctx))
			})
		})

		// Organization
		r.Route("/orgs", func(r chi.Router) {
			handlers.CreateOrganization(r, s.database)
			handlers.ListOrganizations(r)
		})
	})

	s.mux.Group(func(r chi.Router) {
		handlers.Health(s.mux)

		r.Route("/auth/", func(r chi.Router) {
			// Auth
			handlers.Signup(r, s.database)
			handlers.Signin(r, s.database)
		})
	})
}

func (s *Server) requestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		srw := NewStatusResponseWriter(w)

		defer func() {
			s.log.Info(
				"Request sent",
				zap.String("method", req.Method),
				zap.Duration("started at", time.Since(start)),
				zap.Int("status", srw.statusCode),
				zap.String("host", req.Host),
				zap.String("path", req.URL.Path),
				zap.String("query", req.URL.RawQuery),
			)
		}()

		next.ServeHTTP(srw, req)
	})
}
