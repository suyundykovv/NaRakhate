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

func (s *Server) Start(addr string) {
	r := mux.NewRouter()
	logging.Info("Setting up server routes")

	apis := r.PathPrefix("/api").Subrouter()
	apis.Use(api.JWTAuthMiddleware)
	apis.HandleFunc("/protected", api.ProtectedHandler).Methods("GET")

	r.HandleFunc("/getApi", s.FetchFootballMatchesHandler).Methods("GET")
	r.HandleFunc("/addMatch", s.FetchAndSaveMatchesHandler).Methods("GET")
	r.HandleFunc("/addLeague", s.FetchLeagueMatchesHandler).Methods("GET")
	r.HandleFunc("/getMatchs", s.GetAllMatchesHandler).Methods("GET")

	r.HandleFunc("/sign-up", s.SignupHandler).Methods("POST")
	r.HandleFunc("/log-in", s.LoginHandler).Methods("POST")
	r.HandleFunc("/getUser", s.getAllUsersHandler).Methods("GET")
	r.HandleFunc("/deleteUser/{key}", s.deleteUserHandler).Methods("DELETE")
	r.HandleFunc("/updateUser", s.updateUserHandler).Methods("PUT")

	r.HandleFunc("/createBet", s.CreateBetHandler).Methods("POST")
	r.HandleFunc("/getBets", s.GetAllBetsHandler).Methods("GET")
	r.HandleFunc("/getBet/{id}", s.GetBetByIDHandler).Methods("GET")
	r.HandleFunc("/updateBet", s.UpdateBetHandler).Methods("PUT")
	r.HandleFunc("/deleteBet/{id}", s.DeleteBetHandler).Methods("DELETE")

	r.HandleFunc("/createEvent", s.CreateEventHandler).Methods("POST")
	r.HandleFunc("/getEvents", s.GetAllEventsHandler).Methods("GET")
	r.HandleFunc("/getEvent/{id}", s.GetEventByIDHandler).Methods("GET")
	r.HandleFunc("/updateEvent", s.UpdateEventHandler).Methods("PUT")
	r.HandleFunc("/deleteEvent/{id}", s.DeleteEventHandler).Methods("DELETE")

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

func (s *Server) startBackgroundWorker() {
	ticker := time.NewTicker(1 * time.Hour) // Fetch every hour
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			logging.Info("Fetching football matches from Football API Sports...")
			matches, err := s.fetchFootballMatches()
			if err != nil {
				logging.Error("Failed to fetch football match data", err)
				continue
			}

			err = s.saveFootballMatchesToDB(matches)
			if err != nil {
				logging.Error("Failed to save football match data to DB", err)
			}
		case <-s.shutdownCh:
			logging.Info("Background worker shutting down...")
			return
		}
	}
}

func (s *Server) shutdown() {
	close(s.shutdownCh)
}
