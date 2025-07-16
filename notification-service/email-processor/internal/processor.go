package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/hibiken/asynq"
)

const TaskEmail = "notification:email"

type Processor struct{}

func NewProcessor() *Processor {
	return &Processor{}
}

func (p *Processor) ProcessTask(ctx context.Context, task *asynq.Task) error {
	// Simulate 50% failure rate
	if rand.Float64() < 0.5 {
		fmt.Println("email processor failed to process notification")
		return fmt.Errorf("email processor failed to process notification")
	}

	// Unmarshal notification data
	var notification Notification
	if err := json.Unmarshal(task.Payload(), &notification); err != nil {
		return fmt.Errorf("failed to unmarshal notification payload: %w", err)
	}

	// Log the notification (in a real system, would actually send it)
	fmt.Printf("[%s] Successfully processed email notification to %s: %s\n",
		time.Now().Format(time.RFC3339),
		notification.Recipient,
		notification.Content)

	return nil
}
