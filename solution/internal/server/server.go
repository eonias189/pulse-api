package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"solution/internal/db"
)

type Server struct {
	address string
	db      *db.DB
	logger  *slog.Logger
}

func NewServer(address string, db *db.DB, logger *slog.Logger) *Server {
	return &Server{
		address: address,
		logger:  logger,
		db:      db,
	}
}

func (s *Server) LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info(fmt.Sprintf("%v %v", r.Method, r.URL.String()))
		next.ServeHTTP(w, r)
	})
}

func SendError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, err.Error())
}

func (s *Server) handleCountries(w http.ResponseWriter, r *http.Request) {
	region := r.URL.Query().Get("region")
	var (
		countries []db.Country
		err       error
	)
	if region == "" {
		countries, err = s.db.GetCountries()
	} else {
		countries, err = s.db.GetCountriesOfRegion(region)
	}
	if err != nil {
		SendError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	data, _ := json.Marshal(countries)
	fmt.Fprint(w, string(data))
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/ping", s.handlePing)
	mux.HandleFunc("/api/countries", s.handleCountries)

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
