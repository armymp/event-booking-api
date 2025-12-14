package models

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/armymp/event-booking-api/db"
)

// logic for storing and fetching data

type Event struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Location    string    `json:"location" binding:"required"`
	DateTime    time.Time `json:"dateTime" binding:"required"`
	UserID      int       `json:"userId"`
}

func (e *Event) Save() error {
	query := `INSERT INTO events(name, description, location, dateTime, user_id) 
	VALUES (?, ?, ?, ?, ?)`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		slog.Error("Failed to prepare INSERT statement", "error", err)
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserID)
	if err != nil {
		slog.Error("Failed to execute INSERT", "event", e, "error", err)
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		slog.Error("Failed to get LastInsertId", "error", err)
		return err
	}

	e.ID = id
	return nil
}

func GetAllEvents() ([]Event, error) {
	query := "SELECT * FROM events"

	rows, err := db.DB.Query(query)
	if err != nil {
		slog.Error("SELECT query failed", "error", err)
		return nil, err
	}

	defer rows.Close()

	var events []Event

	for rows.Next() {
		var e Event
		err := rows.Scan(&e.ID, &e.Name, &e.Description, &e.Location, &e.DateTime, &e.UserID)
		if err != nil {
			slog.Error("Row scan failed", "error", err)
			return nil, err
		}

		events = append(events, e)
	}

	if err := rows.Err(); err != nil {
		slog.Error("Row iteration error", "error", err)
		return nil, err
	}

	return events, nil
}

func GetEventByID(id int64) (*Event, error) {
	query := "SELECT * FROM events WHERE id = ?"
	row := db.DB.QueryRow(query, id)

	var e Event

	err := row.Scan(&e.ID, &e.Name, &e.Description, &e.Location, &e.DateTime, &e.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Info("Event not found", "event_id", id)
		} else {
			slog.Error("QueryRow failed", "event_id", id, "error", err)
		}
		return nil, err
	}

	return &e, nil
}

func (e Event) Update() error {
	query := `
	UPDATE events
	SET name = ?, description = ?, location = ?, dateTime = ?
	WHERE id = ?
	`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		slog.Error("Failed to prepare UPDATE", "error", err)
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.ID)
	if err != nil {
		slog.Error("Failed to execute UPDATE", "event", e, "error", err)
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		slog.Error("Failed to check RowsAffected", "error", err)
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (e Event) Delete() error {
	query := "DELETE FROM events WHERE id = ?"

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		slog.Error("Failed to prepare DELETE", "error", err)
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(e.ID)
	if err != nil {
		slog.Error("Failed to execute DELETE", "event", e, "error", err)
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		slog.Error("Failed to check RowsAffected", "error", err)
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
