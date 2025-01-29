package api

import (
	"Aitu-Bet/logging"
	"Aitu-Bet/utils"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
)

func ProtectedEndpoint(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Welcome authorized user", "user_id": userID})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.Email == "test@example.com" && request.Password == "password123" {
		token, err := utils.GenerateJWT(1, request.Email)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"token": token})
		return
	}
	http.Error(w, "Invalid credentials", http.StatusUnauthorized)
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to my house body!"))
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
		_, err := s.db.Exec("INSERT INTO data (key, value) VALUES ($1, $2) ON CONFLICT (key) DO UPDATE SET value = $2", key, value)
		if err != nil {
			logging.Error("Failed to insert data into the database", err)
			http.Error(w, "Failed to insert data", http.StatusInternalServerError)
			return
		}
	}
	s.requests++

	logging.Info("New data received and stored", "data", input)
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) getDataHandler(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rows, err := s.db.Query("SELECT key, value FROM data")
	if err != nil {
		logging.Error("Failed to retrieve data from database", err)
		http.Error(w, "Failed to retrieve data", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	data := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			logging.Error("Failed to scan row", err)
			http.Error(w, "Failed to retrieve data", http.StatusInternalServerError)
			return
		}
		data[key] = value
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logging.Error("Failed to encode data into JSON", err)
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

	var value string
	err := s.db.QueryRow("SELECT value FROM data WHERE key = $1", key).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			logging.Warn("Key not found", "key", key)
			http.Error(w, "Key not found", http.StatusNotFound)
		} else {
			logging.Error("Failed to retrieve data", err)
			http.Error(w, "Failed to retrieve data", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{key: value}); err != nil {
		logging.Error("Failed to encode data for specific key", err)
		http.Error(w, "Failed to retrieve data", http.StatusInternalServerError)
		return
	}

	logging.Info("Data for key retrieved", "key", key, "value", value)
}

func (s *Server) deleteDataHandler(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/data/")
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec("DELETE FROM data WHERE key = $1", key)
	if err != nil {
		logging.Error("Failed to delete data from database", err)
		http.Error(w, "Failed to delete data", http.StatusInternalServerError)
		return
	}

	s.requests++

	logging.Info("Data deleted successfully", "key", key)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) statsHandler(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var totalDataEntries int

	err := s.db.QueryRow("SELECT COUNT(*) FROM data").Scan(&totalDataEntries)
	if err != nil {
		logging.Error("Failed to retrieve total data count", err)
		http.Error(w, "Failed to retrieve stats", http.StatusInternalServerError)
		return
	}

	stats := map[string]int{
		"requests": s.requests,
		"size":     totalDataEntries,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		logging.Error("Failed to encode stats into JSON", err)
		http.Error(w, "Failed to retrieve stats", http.StatusInternalServerError)
		return
	}

	logging.Info("Stats successfully retrieved", "stats", stats)
}
