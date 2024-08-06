package utils

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/trebsirk/study-cli/structs"

	"github.com/fatih/color"
)

func GetQuizResponseCLI(q structs.QuizQuestion) int {
	// Prompt the user for an answer
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Your answer (enter the number corresponding to your choice): ")
		input, _ := reader.ReadString('\n')
		input = input[:len(input)-1] // Remove the newline character

		// Validate the input
		answerIndex, err := strconv.Atoi(input)
		if err != nil || answerIndex < 0 || answerIndex >= len(q.CandidateAnswers) {
			fmt.Println("Invalid input. Please enter a number between 0 and", len(q.CandidateAnswers)-1)
			continue
		}

		return answerIndex
	}
}

func AdministerQuizQuestionCLI(q *structs.QuizQuestion) bool {
	fmt.Println("Question:")
	fmt.Println(q.Question)
	fmt.Println("Answer choices:")
	var correctIndex int = -1
	for i, answer := range q.CandidateAnswers {
		fmt.Printf("%d. %s\n", i, answer)
		if answer == q.CorrectAnswer {
			correctIndex = i
		}
	}

	ans := GetQuizResponseCLI(*q)
	if ans == correctIndex {
		// fmt.Println("Correct answer!")
		color.Green("Correct answer!")
		return true
	} else {
		// fmt.Printf("Incorrect answer. The correct answer is: %d. %s\n", correctIndex, q.CorrectAnswer)
		color.Red("Incorrect answer. The correct answer is: %d. %s\n", correctIndex, q.CorrectAnswer)
		return false
	}

}

func SelectQuizQuestionFromDB(db *sql.DB, user string, tags []string) (*structs.QuizQuestion, error) {

	// query := "SELECT id, question, candidate_answers, correct_answer FROM quiz_data WHERE '{\"science\"}' <@ tags LIMIT 1"
	// rows, err := db.Query(query)
	// regarding tags: <@ for AND, && for OR
	// query := "SELECT id, question, candidate_answers, correct_answer FROM quiz_data WHERE $1 <@ tags LIMIT 1"
	// query := "SELECT id, question, candidate_answers, correct_answer FROM quiz_data WHERE tags && $1 LIMIT 1"
	var err error
	var content []byte

	content, err = os.ReadFile("sql/select_next.sql")
	if err != nil {
		log.Fatal("Error reading query file:", err)
	}

	query := string(content)

	if tags == nil {
		//tags = []string{"geography"}
		tagfname := "data/tags.txt"
		fmt.Println("no tags passed ... \n getting tags from file ", tagfname)
		tags, err = ReadFileToList(tagfname)
		fmt.Println("found tags: ", tags)
		//rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
		//i := rand.Intn(len(tags))
		//tags = tags[i : i+1]
		tags = tags[:2] // []string{"S3", "HTML"}
	}
	if err != nil {
		log.Fatal("No tags. Exiting.")
	}
	fmt.Println("using tags: ", tags)
	rows, err := db.Query(query, pq.Array(tags))
	if err != nil {
		log.Fatal("Error querying table quiz_data: ", err)
	}
	defer rows.Close()
	var q structs.QuizQuestion
	var id int
	var question, correct_answer string
	var candidate_questions []string
	for rows.Next() {
		if err := rows.Scan(&id, &question, pq.Array(&candidate_questions), &correct_answer); err != nil {
			log.Println("Error scanning row: ", err)
			return nil, err
		} else {
			q = structs.QuizQuestion{ID: id, Question: question,
				CandidateAnswers: candidate_questions,
				CorrectAnswer:    correct_answer}
			return &q, nil
		}
	}
	return nil, errors.New("query returned null set (no results found for tags)")
}

func InsertQuizResponseToDB(db *sql.DB, q *structs.QuizQuestion, u *structs.User, result bool) error {
	// quiz_id, res, user_id
	content, err := os.ReadFile("sql/insert_into_quiz_results.sql")
	if err != nil {
		log.Fatal("Error reading query file:", err)
	}
	query := string(content)
	qIdStr := strconv.Itoa(q.ID)
	query = strings.ReplaceAll(query, ":quiz_id", qIdStr)
	if result {
		query = strings.ReplaceAll(query, ":res", "true")
	} else {
		query = strings.ReplaceAll(query, ":res", "false")
	}
	uIdStr := strconv.Itoa(u.ID)
	query = strings.ReplaceAll(query, ":user_id", uIdStr)
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal("Error inserting quiz results: ", err)
		return err
	}

	return nil
}
