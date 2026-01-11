package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
}

// init db
// ret db connection and error
func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return db, nil
}

// add a user to the db
func AddUser(db *sql.DB, username, email, password string) (int64, error) {
	if username == "" {
		return 0, fmt.Errorf("username cant be empty")
	}
	if email == "" {
		return 0, fmt.Errorf("email cant be empty")
	}
	if password == "" {
		return 0, fmt.Errorf("password cant be empty")
	}

	// hash passw with bcrypt
	// https://www.slingacademy.com/article/hashing-passwords-securely-with-bcrypt-in-go/
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	// insert user
	insertQuery := `
	INSERT INTO users (username, email, password, created_at)
	VALUES (?, ?, ?, ?)
	`
	result, err := db.Exec(insertQuery, username, email, string(hashedPassword), time.Now())
	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %w", err)
	}

	// get userid
	userID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get user ID: %w", err)
	}

	return userID, nil
}

func Login(db *sql.DB, username, password string) (bool, error) {
	getPwQuery := `
	SELECT password FROM users WHERE username = ?
	`
	var storedPass string
	err := db.QueryRow(getPwQuery, username).Scan(&storedPass)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("user not found")
		}
		return false, fmt.Errorf("failed to get password: %w", err)
	}

	if checkPassword(storedPass, password) {
		return true, nil
	}
	return false, fmt.Errorf("invalid password")
}

func checkPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
