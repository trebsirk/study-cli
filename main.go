package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/briandowns/spinner"
	_ "github.com/fatih/color"
	_ "github.com/lib/pq"
	"github.com/trebsirk/study-cli/structs"
	"github.com/trebsirk/study-cli/utils"
)

func printQuizQuestionFromJSON(fname string) {
	// Open JSON file
	file, err := os.Open(fname)
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

	for _, question := range questions {
		fmt.Println("Question:")
		fmt.Println(question.Question)
		fmt.Println("Answer choices:")
		for i, answer := range question.CandidateAnswers {
			println(i, "-", answer)
		}
		fmt.Println()
	}
}

// to send log output to log file, use go run main.go next 2>> logfile
func main() {
	cmd := os.Args[1]
	fmt.Println("cmd: ", cmd)
	switch cmd {
	case "dev":
		var tags []string = nil
		if len(os.Args) >= 3 {
			tags = os.Args[2:] // []string{}
		}
		fmt.Println("tags: ", tags)
		fmt.Println("fake tags: ", []string{"AWS", "S3"})
	case "load":
		fname := os.Args[2] // eg go run load data/test.json
		utils.LoadFromFile(fname)
	case "printjson":
		fname := os.Args[2]
		printQuizQuestionFromJSON(fname)
	case "next":
		creds, err := utils.ReadCredentialsFromFile()
		if err != nil {
			fmt.Println("error in ReadCredentialsFromFile: ", err)
			return
		}
		var tags []string = nil
		if len(os.Args) >= 3 {
			tags = os.Args[2:]
		}
		config := utils.GetConfig()
		db := utils.GetDB(config)
		defer db.Close()
		q, err := utils.SelectQuizQuestionFromDB(db, creds.Username, tags)
		if err != nil {
			fmt.Println("error selecting quiz question:", err)
			return
		}
		res := utils.AdministerQuizQuestionCLI(q)
		log.Printf("q=%v, r=%t", q.ID, res)
		id, err := utils.GetIdForUsernameFromDB(db, creds.Username)
		if err != nil {
			fmt.Println(err)
			return
		}
		if res {
			idInt, err := strconv.Atoi(id)
			if err != nil {
				fmt.Println(err)
				return
			}
			user := &structs.User{ID: idInt}
			utils.InsertQuizResponseToDB(db, q, user, res)
		}

	case "user":
		// Hash password
		var hashedPassword, err = utils.HashPassword("password1")

		if err != nil {
			println(fmt.Println("Error hashing password"))
			return
		}

		fmt.Println("Password Hash:", hashedPassword)

		// Check if passed password matches the original password
		fmt.Println("Password Match:",
			utils.CheckIfPasswordsMatch(hashedPassword, "password1"))
	case "create-user":
		c, err := utils.ReadCredentialsFromCLI()
		if err != nil {

		}
		// username := "kris"
		// pass := "kris"
		err = utils.CreateUser(&c)
		if err != nil {
			fmt.Println("error creating user:", err)
			return
		}
		okay := utils.ValidateUser(&c)
		if !okay {
			fmt.Println("error validating user", c.Username)
			return
		}
		fmt.Println("user created", c.Username)
	case "stats":
		doneChan := make(chan bool)
		doneAckChan := make(chan bool)
		statsChan := make(chan []structs.Stats)
		timeChan := make(chan time.Time, 2)
		go func() {
			s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
			s.Color("cyan")
			s.Suffix = " getting stats"
			s.Start()
			<-doneChan
			s.Stop()
			// for s.Active() {
			// 	time.Sleep(time.Millisecond * 100)
			// }
			doneAckChan <- true
		}()
		go func() { // get stats
			timeChan <- time.Now()
			stats, err := utils.GetStats()
			timeChan <- time.Now()
			if err != nil {
				fmt.Println("error getting stats:", err)
				return
			}
			time.Sleep(1 * time.Second)
			statsChan <- stats
		}()
		defer close(doneChan)
		defer close(statsChan)
		defer close(doneAckChan)
		select {
		case stats_res := <-statsChan:
			doneChan <- true
			<-doneAckChan
			for _, s := range stats_res {
				fmt.Println(s)
			}
		case <-time.After(3 * time.Second):
			doneChan <- true
			<-doneAckChan
			fmt.Println("get stats timed out")
		}
		a, b := <-timeChan, <-timeChan
		fmt.Println("GetStats took", b.Sub(a))

	case "all-users":
		config := utils.GetConfig()
		db := utils.GetDB(config)
		defer db.Close()
		users, err := utils.GetUsersFromDB(db)
		if err != nil {
			log.Fatal(err)
		}
		for i, user := range users {
			fmt.Println(i, " ", user)
		}
	case "check-creds":
		c, err := utils.ReadCredentialsFromFile()
		if err != nil || c.Password == "" || c.Username == "" {
			valid := false
			for !valid {
				c, err = utils.ReadCredentialsFromCLI()
				if err == nil {
					valid = true
				}
			}
		}
		okay := utils.ValidateUser(&c)
		if !okay {
			fmt.Println("error validating user", c.Username)
			return
		}
		fmt.Println("user", c.Password, "validated")

	}

}
