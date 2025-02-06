package servers

import (
	"Aitu-Bet/internal/models"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func (s *Server) CreateBetHandler(w http.ResponseWriter, r *http.Request) {
	var bet models.Bet
	if err := json.NewDecoder(r.Body).Decode(&bet); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newBet, err := s.createBet(bet)
	if err != nil {
		http.Error(w, "Failed to create bet", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newBet)
}

func (s *Server) GetAllBetsHandler(w http.ResponseWriter, r *http.Request) {
	bets, err := s.readAllBets()
	if err != nil {
		http.Error(w, "Failed to retrieve bets", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bets)
}

func (s *Server) GetBetByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/getBet/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid bet ID", http.StatusBadRequest)
		return
	}

	bet, err := s.readBetByID(id)
	if err != nil {
		http.Error(w, "Bet not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bet)
}

func (s *Server) UpdateBetHandler(w http.ResponseWriter, r *http.Request) {
	var bet models.Bet
	if err := json.NewDecoder(r.Body).Decode(&bet); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := s.updateBet(bet); err != nil {
		http.Error(w, "Failed to update bet", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bet)
}

func (s *Server) DeleteBetHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/deleteBet/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid bet ID", http.StatusBadRequest)
		return
	}

	if err := s.deleteBetData(id); err != nil {
		http.Error(w, "Failed to delete bet", http.StatusInternalServerError)
		return
	}
}
