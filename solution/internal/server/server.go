package server

import (
	"fmt"
	"log/slog"
	"net/http"
)

type Server struct {
	address string
	logger  *slog.Logger
}

func NewServer(address string, logger *slog.Logger) *Server {
	return &Server{
		address: address,
		logger:  logger,
	}
}

func (s *Server) LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info(fmt.Sprintf("%v %v", r.Method, r.URL.String()))
		next.ServeHTTP(w, r)
	})
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/ping", s.handlePing)

	s.logger.Info("server has been started", "address", s.address)

	err := http.ListenAndServe(s.address, s.LogMiddleware(mux))
	if err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) handlePing(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("ok"))
}
