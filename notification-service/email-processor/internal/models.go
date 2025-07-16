package internal

import "time"

// Notification structure that matches the API server's model
type Notification struct {
	Content   string     `json:"content"`
	Channel   string     `json:"channel"`
	TimeZone  string     `json:"timeZone"`
	Recipient string     `json:"recipient"`
	SendAt    *time.Time `json:"sendAt,omitempty"`
}
