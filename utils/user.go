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
	"time"

	"github.com/trebsirk/study-cli/structs"
	"golang.org/x/crypto/bcrypt"
)

var DEFAULT_CREDENTIALS_FILE = ".credentials"

func HashPassword(password string) (string, error) {
	// Convert password string to byte slice
	var passwordBytes = []byte(password)

	// Hash password with bcrypt's min cost
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.MinCost)

	return string(hashedPasswordBytes), err
}

func CheckIfPasswordsMatch(hashedPassword, currPassword string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword), []byte(currPassword))
	return err == nil
}

func GetIdForUsernameFromDB(db *sql.DB, username string) (string, error) {
	stmt, err := db.Prepare("SELECT id FROM users WHERE username = $1")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	// Execute the prepared statement with the provided username
	var id int
	err = stmt.QueryRow(username).Scan(&id)
	if err != nil {
		return "", err
	}
	fmt.Println("found user id ", id)
	idStr := strconv.Itoa(id)
	fmt.Println("after conversion to string: ", idStr)

	return idStr, nil
}

func GetPasswordHashFromDB(db *sql.DB, username string) (string, error) {
	// Prepare the SQL statement with a placeholder for the username
	stmt, err := db.Prepare("SELECT password_hash FROM users WHERE username = $1")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	// Execute the prepared statement with the provided username
	var passwordHash string
	err = stmt.QueryRow(username).Scan(&passwordHash)
	if err != nil {
		return "", err
	}

	return passwordHash, nil
}

func ValidateUser(c *structs.Credentials) bool {
	config := GetConfig()
	db := GetDB(config)
	defer db.Close()
	hashGroundTruth, err := GetPasswordHashFromDB(db, c.Username)
	if err != nil {
		log.Fatal(err)
	}
	return CheckIfPasswordsMatch(hashGroundTruth, c.Password)
}

func CreateUser(c *structs.Credentials) error {
	var hashedPassword, err = HashPassword(c.Password)
	if err != nil {
		fmt.Println("Error hashing password")
		return errors.New("bad password")
	}
	config := GetConfig()
	db := GetDB(config)
	defer db.Close()
	_, err = db.Exec("INSERT INTO users (username, password_hash) VALUES ($1, $2)", c.Username, hashedPassword)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func GetUsersFromDB(db *sql.DB) ([]structs.User, error) {
	// Execute the prepared statement with the provided username
	var id int
	var username, passwordHash string
	var createdAt time.Time
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	users := make([]structs.User, 0)
	for rows.Next() {
		if err := rows.Scan(&id, &username, &passwordHash, &createdAt); err != nil {
			log.Fatal("Error scanning row:", err)
		}
		users = append(users, structs.User{ID: id, Username: username, PasswordHash: passwordHash, CreatedAt: createdAt})
	}

	return users, nil
}

func ReadCredentialsFromFile() (structs.Credentials, error) {
	file, err := os.Open(DEFAULT_CREDENTIALS_FILE)
	if err != nil {
		return structs.Credentials{}, err
	}
	defer file.Close()

	var creds structs.Credentials
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		switch key {
		case "username":
			creds.Username = value
		case "password":
			creds.Password = value
		}
	}
	if err := scanner.Err(); err != nil {
		return structs.Credentials{}, err
	}

	return creds, nil
}

func ReadCredentialsFromCLI() (structs.Credentials, error) {
	// Prompt the user for username
	fmt.Print("Enter username: ")
	reader := bufio.NewReader(os.Stdin)
	username, err := reader.ReadString('\n')
	if err != nil {
		return structs.Credentials{}, err
	}
	username = strings.TrimSpace(username)

	// Prompt the user for password (without showing it on the terminal)
	fmt.Print("Enter password: ")
	password, err := reader.ReadString('\n')
	if err != nil {
		return structs.Credentials{}, err
	}
	password = strings.TrimSpace(password)

	user := structs.Credentials{Username: username, Password: password}

	WriteCredentialsToFile(&user)

	return user, nil
}

func WriteCredentialsToFile(creds *structs.Credentials) error {
	file, err := os.Create(DEFAULT_CREDENTIALS_FILE)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "username=%s\npassword=%s\n", creds.Username, creds.Password)
	if err != nil {
		return err
	}

	fmt.Println("Credentials successfully written to .credentials file.")
	return nil
}

func AcquireUser() {
	c, err := ReadCredentialsFromFile()
	if err != nil || c.Password == "" || c.Username == "" {
		valid := false
		for !valid {
			c, err = ReadCredentialsFromCLI()
			if err == nil {
				valid = true
			}
		}
	}
	okay := ValidateUser(&c)
	if !okay {
		fmt.Println("error validating user", c.Username)
		return
	}
	WriteCredentialsToFile(&c)
}
