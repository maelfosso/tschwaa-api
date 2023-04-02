package handlers

import (
	"net/http"

	"tschwaa.com/api/models"
	"tschwaa.com/api/services"
)

func getCurrentUser(req *http.Request) *models.User {
	user := req.Context().Value(services.JwtUserKey)
	if u, ok := user.(*models.User); ok {
		return u
	}

	return nil
}
