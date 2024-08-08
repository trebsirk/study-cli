package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
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

func InsertSessionToDB(*structs.UserSession) {

}

func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func WriteSessionToFile(sess *structs.UserSession) error {
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
