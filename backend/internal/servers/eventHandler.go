package servers

import (
	"Aitu-Bet/internal/dal"
	"Aitu-Bet/internal/models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func (s *Server) GetEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID := vars["eventId"]
	events, _ := strconv.Atoi(eventID)
	event, err := dal.ReadEventByID(events, s.db)
	if err != nil {
		http.Error(w, "Failed to fetch event", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

func (s *Server) CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	var event models.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newEvent, err := dal.CreateEvent(event, s.db)
	if err != nil {
		http.Error(w, "Failed to create event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newEvent)
}

func (s *Server) GetAllEventsHandler(w http.ResponseWriter, r *http.Request) {
	events, err := dal.ReadAllEvents(s.db)
	if err != nil {
		http.Error(w, "Failed to retrieve events", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func (s *Server) GetEventByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/getEvent/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	event, err := dal.ReadEventByID(id, s.db)
	if err != nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

func (s *Server) UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	var event models.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := dal.UpdateEvent(event, s.db); err != nil {
		http.Error(w, "Failed to update event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(event)
}

func (s *Server) DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/deleteEvent/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	if err := dal.DeleteEventData(id, s.db); err != nil {
		http.Error(w, "Failed to delete event", http.StatusInternalServerError)
		return
	}
}

func (s *Server) GetAllMatchesHandler(w http.ResponseWriter, r *http.Request) {
	fixtures, err := dal.ReadAllFixtures(s.db)
	if err != nil {
		log.Printf("Error reading fixtures: %v", err)
		http.Error(w, "Failed to retrieve fixtures", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(fixtures); err != nil {
		log.Printf("Error encoding fixtures to JSON: %v", err)
		http.Error(w, "Failed to encode fixtures to JSON", http.StatusInternalServerError)
	}
}
