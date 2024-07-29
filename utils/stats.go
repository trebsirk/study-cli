package utils

import (
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/trebsirk/study-cli/structs"
)

func GetStats() ([]structs.Stats, error) {
	config := GetConfig()
	db := GetDB(config)
	defer db.Close()
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
	var date, service string
	var pct float32
	for rows.Next() {
		if err := rows.Scan(&date, &service, &pct); err != nil {
			//return nil, err
			log.Fatal("Error scanning row:", err)
		}
		//fmt.Printf("next topic: %s\n", column1)
		stats = append(stats, structs.Stats{Date: date, Service: service, Pct: pct})
	}
	return stats, nil
}
