package models

import (
	"log/slog"

	"github.com/armymp/event-booking-api/db"
)

type User struct {
	ID       int64
	EMAIL    string `binding:"required"`
	PASSWORD string `binding:"required"`
}

func (u *User) Save() error {
	query := "INSERT INTO users(email, password) VALUES (?,?)"

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		slog.Error("Failed to prepare INSERT statement", "error", err)
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec(u.EMAIL, u.PASSWORD)
	if err != nil {
		slog.Error("failed to execute INSERT statement", "error", err)
		return err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		slog.Error("Failed to get LastInsertId", "error", err)
		return err
	}

	u.ID = userID
	return err
}