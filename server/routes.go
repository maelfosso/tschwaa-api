package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/maelfosso/jwtauth"
	"go.uber.org/zap"
	"tschwaa.com/api/handlers"
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
		handlers.Health(s.mux)

		s.mux.Route("/auth/", func(r chi.Router) {
			// Auth
			handlers.Signup(s.mux, s.database)
			handlers.Signin(s.mux, s.database)
		})
	})

	s.mux.Group(func(r chi.Router) {
		tokenAuth := jwtauth.New("HS512", []byte("schwaa"), nil)
		s.mux.Use(jwtauth.Verifier(tokenAuth))

		s.mux.Use(jwtauth.Authenticator)

		// Organization
		s.mux.Route("/orgs", func(r chi.Router) {
			handlers.CreateOrganization(s.mux, s.database)
			handlers.ListOrganizations(s.mux)
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
