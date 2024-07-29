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
