package servers

import (
	"Aitu-Bet/internal/models"
	"database/sql"
	"errors"
)

func (s *Server) createEvent(event models.Event) (*models.Event, error) {
	err := s.db.QueryRow(
		"INSERT INTO events (name, description, start_time, category) VALUES ($1, $2, $3, $4) RETURNING id",
		event.Name, event.Description, event.StartTime, event.Category).Scan(&event.ID)

	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (s *Server) readAllEvents() ([]models.Event, error) {
	var events []models.Event

	rows, err := s.db.Query("SELECT id, name, description, start_time, category FROM events")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var event models.Event
		if err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.StartTime, &event.Category); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, rows.Err()
}

func (s *Server) readEventByID(id int) (*models.Event, error) {
	var event models.Event

	err := s.db.QueryRow(
		"SELECT id, name, description, start_time, category FROM events WHERE id = $1", id).
		Scan(&event.ID, &event.Name, &event.Description, &event.StartTime, &event.Category)

	if err == sql.ErrNoRows {
		return nil, errors.New("event not found")
	} else if err != nil {
		return nil, err
	}

	return &event, nil
}

func (s *Server) updateEvent(event models.Event) error {
	_, err := s.db.Exec(
		"UPDATE events SET name = $1, description = $2, start_time = $3, category = $4 WHERE id = $5",
		event.Name, event.Description, event.StartTime, event.Category, event.ID)
	return err
}

func (s *Server) deleteEventData(id int) error {
	_, err := s.db.Exec("DELETE FROM events WHERE id = $1", id)
	return err
}
