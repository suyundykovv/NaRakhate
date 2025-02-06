package servers

import (
	"Aitu-Bet/internal/api" // Убедись, что ты правильно импортируешь пакет api
	"Aitu-Bet/logging"
	"context"
	"database/sql"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"os/signal"
	"sync"
)

// Новая структура сервера
type Server struct {
	mu         sync.Mutex
	db         *sql.DB
	requests   int
	shutdownCh chan struct{}
}

// Новый конструктор сервера
func NewServer(db *sql.DB) *Server {
	return &Server{
		db:         db,
		shutdownCh: make(chan struct{}),
	}
}

// Запуск сервера
func (s *Server) Start(addr string) {
	r := mux.NewRouter()
	logging.Info("Setting up server routes")

	// Подключаем маршруты API
	r.HandleFunc("/api/leaderboard", api.GetTopPlayersHandler).Methods("GET")
	r.HandleFunc("/api/player", api.AddPlayerHandler).Methods("POST")

	// Настройка других API маршрутов
	apis := r.PathPrefix("/api").Subrouter()
	apis.Use(api.JWTAuthMiddleware) // Защищаем маршруты
	apis.HandleFunc("/protected", api.ProtectedHandler).Methods("GET")

	// Подключение к серверу
	addr = ":" + addr
	srv := &http.Server{Addr: addr, Handler: r}

	go func() {
		logging.Info("Server is starting", "address", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Error("ListenAndServe failed", err)
		}
	}()

	// Обработка сигнала завершения работы
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

func (s *Server) shutdown() {
	close(s.shutdownCh)
}
