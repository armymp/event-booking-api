package models

import (
	"database/sql"
	"errors"
	"log/slog"

	"github.com/armymp/event-booking-api/db"
	"github.com/armymp/event-booking-api/utils"
)

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"` // '-' ensures password is never returned in responses
}

func (u *User) Save() error {
	query := "INSERT INTO users(email, password) VALUES (?,?)"

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		slog.Error("User.Save: failed to prepare INSERT statement", "error", err)
		return err
	}

	defer stmt.Close()

	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		slog.Error("User.Save: failed to hash password", "error", err)
		return err
	}

	result, err := stmt.Exec(u.Email, hashedPassword)
	if err != nil {
		slog.Error("User.Save: failed to execute INSERT statement", "error", err)
		return err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		slog.Error("User.Save: failed to get last insert id", "error", err)
		return err
	}

	u.ID = userID
	return err
}

func (u *User) ValidateCredentials() error {
	query := "SELECT password FROM users WHERE email = ?"
	row := db.DB.QueryRow(query, u.Email)

	var hashedPassword string
	if err := row.Scan(&hashedPassword); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("invalid credentials")
		}
		
		slog.Error("User.ValidateCredentials: query failed", "error", err)
		return err
	}

	if !utils.CheckPasssword(hashedPassword, u.Password) {
		return errors.New("invalid credentials")
	}

	return nil
}
