package structs

import (
	"time"
)

// QuizQuestion represents a quiz question
type QuizQuestion struct {
	ID               int      `json:"id"`
	Question         string   `json:"question"`
	CandidateAnswers []string `json:"candidate_answers"`
	CorrectAnswer    string   `json:"correct_answer"`
	Tags             []string `json:"tags"`
}

type Stats struct {
	Date    string  `json:"date"`
	Service string  `json:"service"`
	Pct     float32 `json:"pct"`
}

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
}

// Credentials represents the username and password entries
type Credentials struct {
	Username string
	Password string
}

type UserSession struct {
	SessionID int       `db:"session_id" json:"session_id"`
	UserID    int       `db:"user_id" json:"user_id"`
	Token     string    `db:"token" json:"token"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	IPAddress string    `db:"ip_address" json:"ip_address"`
	UserAgent string    `db:"user_agent" json:"user_agent"`
}
