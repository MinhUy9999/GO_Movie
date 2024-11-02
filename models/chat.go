// models/chat.go
package models

import (
	"errors"
	"log"
	"my-app/config"
	"time"
)

// ChatMessage represents a chat message in the database
type ChatMessage struct {
	ChatID      int       `json:"chatID"`
	UserID      int       `json:"userID"`
	MessageText string    `json:"messageText"`
	Timestamp   time.Time `json:"timestamp"`
}

// GetMessages retrieves all chat messages from the database
func GetMessages() ([]ChatMessage, error) {
	rows, err := config.DB.Query("SELECT chatID, userID, messageText, timestamp FROM CHAT ORDER BY timestamp ASC")
	if err != nil {
		log.Printf("Database query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var messages []ChatMessage
	for rows.Next() {
		var msg ChatMessage
		var timestampStr string // Temp variable to scan timestamp as a string

		if err := rows.Scan(&msg.ChatID, &msg.UserID, &msg.MessageText, &timestampStr); err != nil {
			log.Printf("Row scan error: %v", err)
			return nil, err
		}

		// Parse the timestamp string to time.Time
		msg.Timestamp, err = time.Parse("2006-01-02 15:04:05", timestampStr)
		if err != nil {
			log.Printf("Timestamp parsing error: %v", err)
			return nil, err
		}

		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Rows error: %v", err)
		return nil, err
	}

	return messages, nil
}

func AddMessage(userID int, messageText string) error {
	if messageText == "" {
		log.Println("Attempted to save empty message text")
		return errors.New("message text cannot be empty")
	}
	_, err := config.DB.Exec("INSERT INTO CHAT (userID, messageText, timestamp) VALUES (?, ?, ?)", userID, messageText, time.Now())
	if err != nil {
		log.Printf("Database error: %v", err)
	}
	return err
}
