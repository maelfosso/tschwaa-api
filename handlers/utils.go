package handlers

import (
	"net/http"

	"tschwaa.com/api/model"
	"tschwaa.com/api/services"
)

func getCurrentUser(req *http.Request) *model.User {
	user := req.Context().Value(services.JwtUserKey)
	if u, ok := user.(model.User); ok {
		return &u
	}

	return nil
}
