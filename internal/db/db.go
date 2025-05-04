package db

import (
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

// Message represents a stored message in the database
type Message struct {
	ID          string    `db:"id"`
	Content     string    `db:"content"`
	CreatedAt   time.Time `db:"created_at"`
	AccessCount int       `db:"access_count"`
}

// DB is a wrapper around the database connection
type DB struct {
	*sqlx.DB
}

// New creates a new database connection
func New() (*DB, error) {
	// Connect to SQLite database
	db, err := sqlx.Connect("sqlite", "ashnote.db")
	if err != nil {
		return nil, err
	}

	// Create the messages table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS messages (
			id TEXT PRIMARY KEY,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			access_count INTEGER DEFAULT 0
		)
	`)
	if err != nil {
		return nil, err
	}

	log.Println("Database connected and initialized")
	return &DB{db}, nil
}

// StoreMessage stores an encrypted message and returns its unique ID
func (db *DB) StoreMessage(encryptedContent string) (string, error) {
	// Generate a unique ID for this message
	id := uuid.New().String()

	// Store the message in the database
	_, err := db.Exec(
		"INSERT INTO messages (id, content) VALUES (?, ?)",
		id, encryptedContent,
	)
	if err != nil {
		return "", err
	}

	return id, nil
}

// GetMessage retrieves a message by its ID
// If the message has already been accessed, it returns an error
func (db *DB) GetMessage(id string) (string, error) {
	var message Message
	err := db.Get(&message, "SELECT * FROM messages WHERE id = ?", id)
	if err != nil {
		return "", err
	}

	// If the message has already been accessed, return an error
	if message.AccessCount > 0 {
		return "", errors.New("message has already been accessed")
	}

	// Update access count
	_, err = db.Exec("UPDATE messages SET access_count = access_count + 1 WHERE id = ?", id)
	if err != nil {
		return "", err
	}

	return message.Content, nil
}

// DeleteMessage deletes a message from the database
func (db *DB) DeleteMessage(id string) error {
	_, err := db.Exec("DELETE FROM messages WHERE id = ?", id)
	return err
}
