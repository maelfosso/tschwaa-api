package handlers_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/matryer/is"
	"tschwaa.com/api/handlers"
	"tschwaa.com/api/model"
)

type signupperMock struct {
	user model.User
}

func (s *signupperMock) Signup(ctx context.Context, user model.User) error {
	s.user = user
	return nil
}

func TestSignup(t *testing.T) {
	mux := chi.NewMux()
	s := &signupperMock{}
	handlers.Signup(mux, s)

	t.Run("signs up", func(t *testing.T) {
		var jsonBody = `{"firstname":"john","lastname":"doe","email":"john.doe@mail.com","phone":"69032","password":"xwfwe"}`
		is := is.New(t)
		code, resBody := makePostRequest(mux, "/auth/signup", strings.NewReader(jsonBody))
		is.Equal(http.StatusOK, code)
		is.Equal(strings.TrimSpace(resBody), strconv.FormatBool(true))

		decoder := json.NewDecoder(strings.NewReader(jsonBody))
		var user model.User
		if err := decoder.Decode(&user); err != nil {
			panic(err)
		}
		is.Equal(user, s.user)
	})
}

func makePostRequest(handler http.Handler, target string, body io.Reader) (int, string) {
	req := httptest.NewRequest(http.MethodPost, target, body)
	req.Header = createRestHeader()

	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	result := res.Result()
	bodyBytes, err := io.ReadAll(result.Body)
	if err != nil {
		panic(err)
	}
	return result.StatusCode, string(bodyBytes)
}

func createRestHeader() http.Header {
	header := http.Header{}
	header.Set("Content-Type", "application/json")
	return header
}
