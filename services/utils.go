package services

import (
	"os"

	"github.com/maelfosso/jwtauth"
)

type ContextKey struct {
	name string
}

func (k *ContextKey) String() string {
	return "jwtauth context value " + k.name
}

var (
	JwtUserKey                  = &ContextKey{"JWT_USER"}
	TokenAuth  *jwtauth.JWTAuth = jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), nil)
)
