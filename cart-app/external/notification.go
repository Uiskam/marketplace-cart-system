package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type NotificationClient struct {
	baseURL string
	client  *http.Client
}

type Notification struct {
	Content   string `json:"content" binding:"required"`
	Channel   string `json:"channel" binding:"required,oneof=email push"`
	Recipient string `json:"recipient" binding:"required"`
	//SendAt    *time.Time `json:"sendAt,omitempty" json:"-"`
	Priority string `json:"priority" binding:"required,oneof=high low"`
}

func NewEmailNotification(content, recipient string) *Notification {
	return &Notification{
		Content:   content,
		Channel:   "email",
		Recipient: recipient,
		Priority:  "low",
	}
}

func NewNotificationClient(baseURL string) *NotificationClient {
	return &NotificationClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (nc *NotificationClient) SendTask(taskType string, notification *Notification) error {
	jsonData, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal task request: %w", err)
	}

	url := fmt.Sprintf("%s/tasks", nc.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := nc.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("notification service returned status: %d", resp.StatusCode)
	}

	return nil
}
