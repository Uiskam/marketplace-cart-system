package internal

import (
	"github.com/hibiken/asynq"
	"time"
)

type Notification struct {
	Content   string     `json:"content" binding:"required"`
	Channel   string     `json:"channel" binding:"required,oneof=email push"`
	Recipient string     `json:"recipient" binding:"required"`
	SendAt    *time.Time `json:"sendAt,omitempty" json:"-"`
	Priority  string     `json:"priority" binding:"required,oneof=high low"`
}

// TaskInfo represents task information for the API responses
type TaskInfo struct {
	ID          string        `json:"id"`
	Type        string        `json:"type"`
	Payload     string        `json:"payload"`
	Queue       string        `json:"queue"`
	Retried     int           `json:"retried"`
	Timeout     time.Duration `json:"timeout"`
	State       string        `json:"state"`
	CompletedAt time.Time     `json:"completed_at,omitempty"`
}

func NewTaskInfo(taskInfo *asynq.TaskInfo) *TaskInfo {
	t := &TaskInfo{}
	t.ID = taskInfo.ID
	t.Type = taskInfo.Type
	t.Payload = string(taskInfo.Payload)
	t.Queue = taskInfo.Queue
	t.Retried = taskInfo.Retried
	t.Timeout = taskInfo.Timeout
	t.State = taskInfo.State.String()
	t.CompletedAt = taskInfo.CompletedAt
	return t
}

// TaskListResponse is the response for pending/completed/failed task endpoints
type TaskListResponse struct {
	Tasks []*TaskInfo `json:"tasks"`
	Count int         `json:"count"`
}
