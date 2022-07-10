package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type Server struct {
	address string
	log     *zap.Logger
	mux     chi.Router
	server  *http.Server
}

type Options struct {
	Host string
	Log  *zap.Logger
	Port int
}
