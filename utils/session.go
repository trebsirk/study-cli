package utils

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/trebsirk/study-cli/structs"
)

var DEFAULT_SESSION_FILE = ".session"
var TOKEN_LEN = 32

func CreateUserSession(id int) *structs.UserSession {
	token := GenerateSecureToken(TOKEN_LEN)
	created_at := time.Now()
	expires_at := created_at.Add(time.Hour * 24 * 7)
	return &structs.UserSession{UserID: id, Token: token, CreatedAt: created_at, ExpiresAt: expires_at}
}

func InsertSessionToDB(db *sql.DB, sess *structs.UserSession) error {
	var err error
	var content []byte

	content, err = os.ReadFile("sql/insert_session_info.sql")
	if err != nil {
		log.Fatal("Error reading query file:", err)
	}
	query := string(content)

	query = strings.ReplaceAll(query, ":user_id", strconv.Itoa(sess.UserID))
	query = strings.ReplaceAll(query, ":token", sess.Token)
	query = strings.ReplaceAll(query, ":created_at", sess.CreatedAt.Format("2006-01-02 15:04:05"))
	query = strings.ReplaceAll(query, ":expires_at", sess.ExpiresAt.Format("2006-01-02 15:04:05"))

	_, err = db.Exec(query)
	if err != nil {
		log.Fatal("Error inserting session info: ", err)
		return err
	}

	return nil
}

func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func WriteSessionToFile(sess *structs.UserSession) error {
	// write session info to DEFAULT_SESSION_FILE
	file, err := os.Create(DEFAULT_SESSION_FILE)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "token=%s\ncreated_at=%s\nexpires_at=%s\n",
		sess.Token,
		sess.CreatedAt,
		sess.ExpiresAt)
	if err != nil {
		return err
	}

	fmt.Println("Credentials successfully written to .credentials")
	return nil

}
