package servers

import (


	"Aitu-Bet/internal/api"
	"Aitu-Bet/logging"
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	maxRetries = 5
	retryDelay = 5 * time.Second
)

type Server struct {
	mu         sync.Mutex
	db         *sql.DB
	requests   int
	shutdownCh chan struct{}
}

func NewServer(db *sql.DB) *Server {
	return &Server{
		db:         db,
		shutdownCh: make(chan struct{}),
	}
}

func (s *Server) startBackgroundWorker() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.mu.Lock()
			logging.Info("Server status", "requests", s.requests)
			s.mu.Unlock()
		case <-s.shutdownCh:
			logging.Info("Background worker shutting down...")
			return
		}
	}
}

func (s *Server) shutdown() {
	close(s.shutdownCh)
}

func (s *Server) Start(addr string) {
	r := mux.NewRouter()
	logging.Info("Setting up server routes")

	logging.Info("Received request to start server", "address", addr)
	r.HandleFunc("/login", api.LoginHandler).Methods("POST")

	apis := r.PathPrefix("/api").Subrouter()
	apis.Use(api.JWTAuthMiddleware)
	apis.HandleFunc("/protected", api.ProtectedHandler).Methods("GET")

	r.HandleFunc("/data", s.api.postDataHandler).Methods("POST")
	r.HandleFunc("/data", api.s.getDataHandler).Methods("GET")
	r.HandleFunc("/data/{key}", s.getDatasHandler).Methods("GET")
	r.HandleFunc("/data/{key}", s.deleteDataHandler).Methods("DELETE")
	r.HandleFunc("/stats", s.statsHandler).Methods("GET")
	r.HandleFunc("/", s.getDataHandler).Methods("GET")

	go s.startBackgroundWorker()

	addr = ":" + addr
	srv := &http.Server{Addr: addr, Handler: r}

	go func() {
		logging.Info("Server is starting", "address", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Error("ListenAndServe failed", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	logging.Info("Shutting down server gracefully...")
	s.shutdown()
	if err := srv.Shutdown(context.Background()); err != nil {
		logging.Error("Server shutdown failed", err)
	} else {
		logging.Info("Server gracefully shut down")
	}
}

