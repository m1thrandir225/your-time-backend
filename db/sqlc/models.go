// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID           uuid.UUID      `json:"id"`
	Title        string         `json:"title"`
	DueDate      time.Time      `json:"due_date"`
	ReminderDate sql.NullTime   `json:"reminder_date"`
	Description  sql.NullString `json:"description"`
	UserID       uuid.UUID      `json:"user_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
