package servers

import (
	"Aitu-Bet/logging"
	"Aitu-Bet/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	mu         sync.Mutex
	data       map[string]string
	requests   int
	shutdownCh chan struct{}
}

func NewServer() *Server {
	return &Server{
		data:       make(map[string]string),
		shutdownCh: make(chan struct{}),
	}
}

func (s *Server) postDataHandler(w http.ResponseWriter, r *http.Request) {
	var input map[string]string
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logging.Error("Failed to decode JSON", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for key, value := range input {
		s.data[key] = value
	}
	s.requests++

	logging.Info("Received new data", "data", input)
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) getDataHandler(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(s.data); err != nil {
		logging.Error("Failed to encode data", err)
		http.Error(w, "Failed to retrieve data", http.StatusInternalServerError)
		return
	}

	logging.Info("Data successfully retrieved")
}

func (s *Server) getDatasHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	s.mu.Lock()
	defer s.mu.Unlock()

	value, exists := s.data[key]
	if !exists {
		logging.Warn("Key not found", "key", key)
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{key: value}); err != nil {
		logging.Error("Failed to encode specific data", err)
		http.Error(w, "Failed to retrieve data", http.StatusInternalServerError)
		return
	}

	logging.Info("Data for key retrieved", "key", key, "value", value)
}

func (s *Server) statsHandler(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	stats := map[string]int{
		"requests": s.requests,
		"size":     len(s.data),
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		logging.Error("Failed to encode stats", err)
		http.Error(w, "Failed to retrieve stats", http.StatusInternalServerError)
		return
	}

	logging.Info("Stats successfully retrieved", "stats", stats)
}

func (s *Server) deleteDataHandler(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/data/")
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[key]; !exists {
		logging.Warn("Attempted to delete non-existent key", "key", key)
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	delete(s.data, key)
	s.requests++

	logging.Info("Data successfully deleted", "key", key)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) startBackgroundWorker() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.mu.Lock()
			logging.Info("Server status", "requests", s.requests, "data_entries", len(s.data))
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
	logging.Info("Received request to start server on path", "path", r.Path)
	r.HandleFunc("/data", s.postDataHandler).Methods("POST")
	r.HandleFunc("/data", s.getDataHandler).Methods("GET")
	r.HandleFunc("/data/{key}", s.getDatasHandler).Methods("GET")
	r.HandleFunc("/data/{key}", s.deleteDataHandler).Methods("DELETE")
	r.HandleFunc("/stats", s.statsHandler).Methods("GET")

	go s.startBackgroundWorker()

	addr = ":" + addr
	srv := &http.Server{Addr: addr}

	// Start server in a goroutine
	go func() {
		logging.Info("Server is starting", "address", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Error("ListenAndServe failed", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	logging.Info("Shutting down server...")
	s.shutdown()
	if err := srv.Shutdown(context.Background()); err != nil {
		logging.Error("Server shutdown failed", err)
	} else {
		logging.Info("Server gracefully shut down")
	}
}

func splitPath(path string) (string, string, error) {
	defer utils.CatchCriticalPoint()

	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		return "", "", fmt.Errorf("invalid path: %s", path)
	}

	item := parts[1]
	itemId := parts[2]

	return item, itemId, nil
}
