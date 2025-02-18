package servers

import (
	"Aitu-Bet/internal/dal"
	"Aitu-Bet/internal/models"
	"Aitu-Bet/internal/services"
	"Aitu-Bet/logging"
	"encoding/json"
	"log/slog"
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

	if bet.EventID == 0 {
		http.Error(w, "Event ID is required", http.StatusBadRequest)
		return
	}
	if bet.OddSelection != "home" && bet.OddSelection != "away" && bet.OddSelection != "draw" {
		http.Error(w, "Invalid odd selection; must be 'home', 'away', or 'draw'", http.StatusBadRequest)
		return
	}
	if bet.OddValue <= 0 {
		http.Error(w, "Odd value must be greater than 0", http.StatusBadRequest)
		return
	}

	if err := services.CheckEventStatusForBet(bet.EventID, s.db); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logging.Info("Fetching football matches from Football API Sports...")
	matches, err := s.fetchFootballMatches()
	if err != nil {
		logging.Error("Failed to fetch football match data", err)
		return
	}

	err = s.saveFootballMatchesToDB(matches)
	if err != nil {
		logging.Error("Failed to save football match data to DB", err)
	}

	newBet, err := dal.CreateBet(bet, s.db)
	if err != nil {
		http.Error(w, "Failed to create bet", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newBet)
}

func (s *Server) GetAllBetsHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("start of popa1")

	bets, err := dal.ReadAllBets(s.db)
	if err != nil {
		http.Error(w, "Failed to retrieve bets", http.StatusInternalServerError)
		return
	}
	slog.Info("start of popa1")

	logging.Info("Fetching football matches from Football API Sports...")
	matches, err := s.fetchFootballMatches()
	if err != nil {
		logging.Error("Failed to fetch football match data", err)
		return
	}
	slog.Info("start of popa1")

	err = s.saveFootballMatchesToDB(matches)
	if err != nil {
		logging.Error("Failed to save football match data to DB", err)
	}
	slog.Info("start of popa1")

	err = services.UpdateAllBetsIfMatchFinished(s.db)
	if err != nil {
		logging.Error("Failed to update all bets status match data to DB", err)
	}
	slog.Info("start of popa1")
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
	logging.Info("Fetching football matches from Football API Sports...")
	matches, err := s.fetchFootballMatches()
	if err != nil {
		logging.Error("Failed to fetch football match data", err)
		return
	}

	err = s.saveFootballMatchesToDB(matches)
	if err != nil {
		logging.Error("Failed to save football match data to DB", err)
	}
	bet, err := dal.ReadBetByID(id, s.db)
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

	if err := dal.UpdateBet(bet, s.db); err != nil {
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

	if err := dal.DeleteBetData(id, s.db); err != nil {
		http.Error(w, "Failed to delete bet", http.StatusInternalServerError)
		return
	}
}
