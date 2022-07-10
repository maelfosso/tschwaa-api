package server

import (
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
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

func New(opts Options) *Server {
	if opts.Log == nil {
		opts.Log = zap.NewNop()
	}

	address := net.JoinHostPort(opts.Host, strconv.Itoa(opts.Port))
	mux := chi.NewMux()

	return &Server{
		address: address,
		log:     opts.Log,
		mux:     mux,
		server: &http.Server{
			Addr:              address,
			Handler:           mux,
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      5 * time.Second,
			IdleTimeout:       5 * time.Second,
		},
	}
}
