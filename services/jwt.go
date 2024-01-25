package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type contextKey struct {
	name string
}

var JWTMemberKey *contextKey
var JWTClaimsKey *contextKey
var JWTTokenKey *contextKey
var JWTErrorKey *contextKey
var jwtSecretKey []byte // (os.Getenv("JWT_SECRET"))

type TWAJWTClaims struct {
	*jwt.RegisteredClaims
	User interface{}
}

func (k *contextKey) String() string {
	return "jwtauth context value " + k.name
}

func init() {
	JWTMemberKey = &contextKey{"Member"}
	JWTClaimsKey = &contextKey{"Claims"}
	JWTTokenKey = &contextKey{"Token"}
	JWTErrorKey = &contextKey{"Error"}
	jwtSecretKey = []byte(os.Getenv("JWT_SECRET"))
	// TokenAuth = jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), "s-tschwaa")
}

func GenerateJWTToken(data map[string]interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	now := time.Now().UTC()
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = now.Add(60 * time.Hour).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()
	// claims["authorized"] = true
	claims["user"] = data

	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func TokenFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		if err == http.ErrNoCookie {
			return "", fmt.Errorf("no cookie found")
		} else {
			return "", fmt.Errorf("error when getting cookies: %w", err)
		}
	}
	if cookie.Value == "" {
		return "", fmt.Errorf("empty token found")
	}

	return cookie.Value, nil
}

func TokenFromHeader(r *http.Request) (string, error) {
	var tokenString string

	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		tokenString = bearer[7:]
	} else {
		return "", fmt.Errorf("no authorization found in header")
	}

	if tokenString == "" {
		return "", fmt.Errorf("empty token found")
	}

	return tokenString, nil
}

func Verifier(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := TokenFromCookie(r)

		ctx := r.Context()
		ctx = context.WithValue(ctx, JWTTokenKey, tokenString)
		ctx = context.WithValue(ctx, JWTErrorKey, err)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ParseJWTToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tokenString, _ := ctx.Value(JWTTokenKey).(string)
		err, _ := ctx.Value(JWTErrorKey).(error)

		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return jwtSecretKey, nil
		})

		if err != nil {
			ctx = context.WithValue(ctx, JWTErrorKey, fmt.Errorf("invalidate token: %v", err))
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		if token == nil || !token.Valid {
			ctx = context.WithValue(ctx, JWTErrorKey, errors.New("invalid token"))
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx = context.WithValue(ctx, JWTErrorKey, errors.New("invalid token claims"))
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		ctx = context.WithValue(ctx, JWTClaimsKey, claims["user"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err, _ := ctx.Value(JWTErrorKey).(error)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
