package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type contextKey struct {
	name string
}

var JwtUserKey *contextKey
var JwtClaimsKey *contextKey
var JwtTokenKey *contextKey
var JwtErrorKey *contextKey
var jwtSecretKey []byte // (os.Getenv("JWT_SECRET"))

// var TokenAuth *jwtauth.JWTAuth //  = jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), nil)

func (k *contextKey) String() string {
	return "jwtauth context value " + k.name
}

func init() {
	JwtUserKey = &contextKey{"User"}
	JwtClaimsKey = &contextKey{"Claims"}
	JwtTokenKey = &contextKey{"Token"}
	JwtErrorKey = &contextKey{"Error"}
	jwtSecretKey = []byte(os.Getenv("JWT_SECRET"))
	// TokenAuth = jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), "s-tschwaa")
}

func GenerateJWTToken(data map[string]interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	now := time.Now().UTC()
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = now.Add(60 * time.Minute).Unix()
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

func extractTokenFromRequest(r *http.Request) (string, error) {
	var tokenString string

	authorizationHeader := r.Header.Get("Authorization")
	log.Println("Authorization Bearer : ", authorizationHeader)
	if strings.HasPrefix(authorizationHeader, "Bearer ") {
		tokenString = strings.TrimPrefix(authorizationHeader, "Bearer ")
	} else {
		return "", fmt.Errorf("no authorization found in header")
	}

	if tokenString == "" {
		return "", fmt.Errorf("empty authorization found")
	}

	return tokenString, nil
}

func parseJWTToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return jwtSecretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalidate token: %v", err)
	}

	if token == nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// if token.Valid {
	// 	fmt.Println("You look nice today")
	// } else if errors.Is(err, jwt.ErrTokenMalformed) {
	// 	fmt.Println("That's not even a token")
	// } else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
	// 	// Token is either expired or not active yet
	// 	fmt.Println("Timing is everything")
	// } else {
	// 	fmt.Println("Couldn't handle this token:", err)
	// }

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// exp := claims["exp"].(float64)
	// if int64(exp) < time.Now().Local().Unix() {
	// 	return nil, errors.New("token expired")
	// }

	return &claims, nil
}

func Verifier(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := extractTokenFromRequest(r)

		ctx := r.Context()
		ctx = context.WithValue(ctx, JwtTokenKey, tokenString)
		ctx = context.WithValue(ctx, JwtErrorKey, err)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token, _ := ctx.Value(JwtTokenKey).(string)
		err, _ := ctx.Value(JwtErrorKey).(error)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		claims, err := parseJWTToken(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, JwtClaimsKey, (*claims)["user"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
