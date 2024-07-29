package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/trebsirk/study-cli/structs"
)

func insertQuizQuestion(db *sql.DB, q structs.QuizQuestion) error {
	res, err := db.Exec(
		`INSERT INTO quiz_data (question, candidate_answers, correct_answer, tags) 
		VALUES ($1, $2, $3, $4) 
		ON CONFLICT (question) DO NOTHING;`,
		q.Question, pq.Array(q.CandidateAnswers), q.CorrectAnswer, pq.Array(q.Tags))
	if err != nil {
		fmt.Println("db.Exec result:", res)
		return err
	}

	return nil
}

func LoadFromFile(filename string) {

	// Open JSON file
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read JSON file
	jsonData, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	// Parse JSON data
	var questions []structs.QuizQuestion
	if err := json.Unmarshal(jsonData, &questions); err != nil {
		log.Fatal(err)
	}

	config := GetConfig()
	db := GetDB(config)
	defer db.Close()

	// Insert quiz questions into the database
	for _, question := range questions {
		err := insertQuizQuestion(db, question)
		if err != nil {
			fmt.Println("err:", err)
		}
	}

	fmt.Println("Quiz questions inserted successfully!")

}
