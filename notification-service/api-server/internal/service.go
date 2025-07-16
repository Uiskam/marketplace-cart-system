package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hibiken/asynq"
	"time"
)

const ServiceContextKey = "service"

type Service struct {
	client    *asynq.Client
	inspector *asynq.Inspector
}

const (
	TaskEmail = "notification:email"
	TaskPush  = "notification:push"
)

func NewService(client *asynq.Client, inspector *asynq.Inspector) *Service {
	return &Service{
		client:    client,
		inspector: inspector,
	}
}

// EnqueueNotification enqueues a notification task for processing.
func (s *Service) EnqueueNotification(notification *Notification) error {
	opts := []asynq.Option{
		asynq.MaxRetry(3),
		asynq.Queue(notification.Channel + "_" + notification.Priority),
		asynq.Retention(time.Hour * 24 * 365),
	}

	if notification.SendAt != nil {
		hour := notification.SendAt.Hour()
		if hour >= 22 || hour < 6 {
			adjusted := time.Date(
				notification.SendAt.Year(),
				notification.SendAt.Month(),
				notification.SendAt.Day(),
				6, 0, 0, 0,
				notification.SendAt.Location())
			if hour >= 22 {
				adjusted = adjusted.AddDate(0, 0, 1)
			}
			notification.SendAt = &adjusted
		}
		opts = append(opts, asynq.ProcessAt(*notification.SendAt))
	}
	payload, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to serialize notification: %w", err)
	}
	taskType := TaskEmail
	if notification.Channel == "push" {
		taskType = TaskPush
	}
	task := asynq.NewTask(taskType, payload)
	_, err = s.client.Enqueue(task, opts...)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	return nil
}

// GetPendingTasks retrieves all tasks that are pending or scheduled for future processing.
func (s *Service) GetPendingTasks() (*TaskListResponse, error) {
	queues, err := s.inspector.Queues()
	if err != nil {
		return nil, fmt.Errorf("failed to get queues: %w", err)
	}

	var tasks []*TaskInfo
	totalCount := 0

	for _, queue := range queues {
		pendingTasks, err := s.inspector.ListPendingTasks(queue, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to list pending tasks: %w", err)
		}
		for _, task := range pendingTasks {
			tasks = append(tasks, NewTaskInfo(task))
		}

		scheduledTasks, err := s.inspector.ListScheduledTasks(queue, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to list scheduled tasks: %w", err)
		}
		for _, task := range scheduledTasks {
			tasks = append(tasks, NewTaskInfo(task))
		}

		retryTasks, err := s.inspector.ListRetryTasks(queue, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to list retry tasks: %w", err)
		}
		for _, task := range retryTasks {
			tasks = append(tasks, NewTaskInfo(task))
		}

		totalCount += len(pendingTasks) + len(scheduledTasks)
	}

	return &TaskListResponse{
		Tasks: tasks,
		Count: totalCount,
	}, nil
}

// GetCompletedTasks retrieves tasks that have been successfully processed.
func (s *Service) GetCompletedTasks() (*TaskListResponse, error) {
	queues, err := s.inspector.Queues()
	if err != nil {
		return nil, fmt.Errorf("failed to get queues: %w", err)
	}

	var tasks []*TaskInfo
	for _, queue := range queues {
		completedTasks, err := s.inspector.ListCompletedTasks(queue)
		if err != nil {
			return nil, fmt.Errorf("failed to list completed tasks: %w", err)
		}
		for _, task := range completedTasks {
			tasks = append(tasks, NewTaskInfo(task))
		}

	}
	return &TaskListResponse{
		Tasks: tasks,
		Count: len(tasks),
	}, nil
}

// GetFailedTasks retrieves tasks that have failed all retry attempts.
func (s *Service) GetFailedTasks() (*TaskListResponse, error) {
	queues, err := s.inspector.Queues()
	if err != nil {
		return nil, fmt.Errorf("failed to get queues: %w", err)
	}

	var tasks []*TaskInfo
	for _, queue := range queues {
		failedTasks, err := s.inspector.ListArchivedTasks(queue)
		if err != nil {
			return nil, fmt.Errorf("failed to list failed tasks: %w", err)
		}
		for _, task := range failedTasks {
			tasks = append(tasks, NewTaskInfo(task))
		}

	}
	return &TaskListResponse{
		Tasks: tasks,
		Count: len(tasks),
	}, nil
}

// SendNow sends a task immediately.
func (s *Service) SendNow(taskID string) (bool, error) {
	queues, err := s.inspector.Queues()
	if err != nil {
		return false, fmt.Errorf("failed to get queues: %w", err)
	}

	for _, queue := range queues {
		err = s.inspector.RunTask(queue, taskID)
		if err != nil {
			if errors.Is(err, asynq.ErrTaskNotFound) {
				continue
			}
			return false, fmt.Errorf("failed to run task: %w", err)
		}
		return true, nil
	}

	return false, nil
}

// CancelTask cancels a task.
func (s *Service) CancelTask(taskID string) (bool, error) {
	queues, err := s.inspector.Queues()
	if err != nil {
		return false, fmt.Errorf("failed to get queues: %w", err)
	}

	for _, queue := range queues {
		err = s.inspector.ArchiveTask(queue, taskID)
		if err != nil {
			if errors.Is(err, asynq.ErrTaskNotFound) {
				continue
			}
			return false, fmt.Errorf("failed to cancel task: %w", err)
		}
		return true, nil
	}

	return false, nil
}
