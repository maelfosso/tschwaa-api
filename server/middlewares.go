package server

import (
	"context"
	"net/http"
	"time"

	"github.com/maelfosso/jwtauth"
	"go.uber.org/zap"
	"tschwaa.com/api/services"
)

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

func (s *Server) convertJwtTokenToUser(next http.Handler) http.Handler {
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
}
