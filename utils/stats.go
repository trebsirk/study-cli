package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/trebsirk/study-cli/structs"
)

func GetStats(db *sql.DB) ([]structs.Stats, error) {
	// query := "SELECT id, question, candidate_answers, correct_answer FROM quiz_data WHERE '{\"science\"}' <@ tags LIMIT 1"
	// rows, err := db.Query(query)
	// regarding tags: <@ for AND, && for OR
	// query := "SELECT id, question, candidate_answers, correct_answer FROM quiz_data WHERE $1 <@ tags LIMIT 1"
	// query := "SELECT id, question, candidate_answers, correct_answer FROM quiz_data WHERE tags && $1 LIMIT 1"
	var err error
	var content []byte

	content, err = os.ReadFile("sql/stats.sql")
	if err != nil {
		log.Fatal("Error reading query file:", err)
	}
	query := string(content)

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal("Error querying table quiz_data:", err)
	}
	defer rows.Close()
	stats := make([]structs.Stats, 0)
	var stat = structs.Stats{}
	for rows.Next() {
		if err := rows.Scan(&stat.Date, &stat.Service, &stat.Pct); err != nil {
			//return nil, err
			log.Fatal("Error scanning row:", err)
		}
		//fmt.Printf("next topic: %s\n", column1)
		stats = append(stats, stat) //structs.Stats{Date: date, Service: service, Pct: pct})
	}
	return stats, nil
}

// WriteStatsToFile writes an array of Stats structs to a JSON file.
func WriteStatsToFile(stats []structs.Stats, filename string) error {
	// Open the file for writing
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Create a JSON encoder and write the data to the file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty print with indentation

	if err := encoder.Encode(stats); err != nil {
		return fmt.Errorf("failed to encode stats to JSON: %w", err)
	}

	return nil
}
