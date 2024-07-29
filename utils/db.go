package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func (c *Config) PrintValues() {
	fmt.Println("Host:", c.Host)
	fmt.Println("Port:", c.Port)
	fmt.Println("User:", c.User)
	fmt.Println("Password:", "***")
	fmt.Println("DBName:", c.DBName)
}

func GetConfig() *Config {
	return &Config{
		Host:     os.Getenv("PG_HOST"),
		Port:     os.Getenv("PG_PORT"),
		User:     os.Getenv("PG_USER"),
		Password: os.Getenv("PG_PASSWORD"),
		DBName:   os.Getenv("PG_DBNAME"),
	}
}

func GetDB(c *Config) *sql.DB {
	// Construct the connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Password, c.DBName)

	// Connect to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	return db
}
