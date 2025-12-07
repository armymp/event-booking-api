package models

import (
	"time"

	"github.com/armymp/event-booking-api/db"
)

// logic for storing and fetching data

type Event struct {
	ID          int64       `json:"id"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Location    string    `json:"location" binding:"required"`
	DateTime    time.Time `json:"dateTime" binding:"required"`
	UserID      int       `json:"userId"`
}

var events = []Event{}

func (e Event) Save() error {
	query := `INSERT INTO events(name, description, location, dateTime, user_id) 
	VALUES (?, ?, ?, ?, ?)`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserID)

	if err != nil {
		return err
	}
	
	id, err := result.LastInsertId()
	e.ID = id

	return err
}

func GetAllEvents() []Event {
	return events
}
