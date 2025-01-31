package servers

import (
	"Aitu-Bet/internal/api"
	"Aitu-Bet/internal/models"
	"Aitu-Bet/logging"
	"Aitu-Bet/utils"
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
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

	r.HandleFunc("/sign-up", s.SignupHandler).Methods("POST")
	r.HandleFunc("/log-in", s.LoginHandler).Methods("POST")
	r.HandleFunc("/getUser", s.getAllUsersHandler).Methods("GET")
	r.HandleFunc("/deleteUser/{key}", s.deleteUserHandler).Methods("DELETE")

	r.HandleFunc("/data", s.postDataHandler).Methods("POST")
	r.HandleFunc("/data", s.getDataHandler).Methods("GET")
	r.HandleFunc("/data/{key}", s.getDatasHandler).Methods("GET")
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
func (s *Server) getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := s.GetAllUsers()
	if err != nil {
		logging.Error("Failed to retrieve users", err)
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}

	// Return the list of users as a JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		logging.Error("Failed to encode users into JSON", err)
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}

	logging.Info("Users successfully retrieved")
}
func (s *Server) deleteUserData(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	log.Println("this is key", key)
	_, err := s.db.Exec("DELETE FROM users WHERE key = $1", key)
	if err != nil {
		logging.Error("Failed to delete user from database", err, key)
		return err
	}

	s.requests++

	logging.Info("Data deleted successfully", "key", key)
	return err
}
func (s *Server) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/deleteUser/")
	err := s.deleteUserData(key)
	if err != nil {
		logging.Error("Failed to delete users", err)
		http.Error(w, "Failed to delete users", http.StatusInternalServerError)
		return
	}

	logging.Info("Users successfully deleted")
}
func (s *Server) GetAllUsers() ([]models.User, error) {
	var users []models.User

	rows, err := s.db.Query("SELECT id, username, email, role FROM users")
	if err != nil {
		logging.Error("Failed to query users from database", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role); err != nil {
			logging.Error("Failed to scan row into user struct", err)
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		logging.Error("Error iterating over rows", err)
		return nil, err
	}

	return users, nil
}

func (s *Server) SignupHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		logging.Error("Failed to hash password", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var user models.User
	err = s.db.QueryRow("INSERT INTO users (username, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id, username, email, role",
		request.Username, request.Email, hashedPassword, request.Role).Scan(&user.ID, &user.Username, &user.Email, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Error inserting new user", http.StatusInternalServerError)
			return
		}
		logging.Error("Failed to insert user into database", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var user models.User
	err := s.db.QueryRow("SELECT id, username, email, password, role FROM users WHERE email = $1", request.Email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		logging.Error("Failed to retrieve user from database", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1))

		_, err := utils.ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
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
