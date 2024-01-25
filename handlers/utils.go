package handlers

import (
	"net/http"

	"tschwaa.com/api/models"
	"tschwaa.com/api/services"
)

func GetCurrentMember(req *http.Request) *models.Member {
	user := req.Context().Value(services.JWTMemberKey)
	if user == nil {
		return nil
	}

	if u, ok := user.(*models.Member); ok {
		return u
	}

	return nil
}
