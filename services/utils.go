package services

import (
	"os"

	"github.com/maelfosso/jwtauth"
)

var TokenAuth *jwtauth.JWTAuth = jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), nil)
