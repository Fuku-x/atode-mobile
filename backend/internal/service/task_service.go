package service

import (
	"context"
	"errors"
	"strings"

	"example.com/atode/backend/internal/model"
	"github.com/google/uuid"
)

var ErrInvalidTaskTitle = errors.New("invalid task title")

type TaskCreator interface {
	Create(ctx context.Context, userID uuid.UUID, title string) (model.Task, error)
}

type TaskService struct {
	repo TaskCreator
}

func NewTaskService(repo TaskCreator) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(ctx context.Context, userID uuid.UUID, title string) (model.Task, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return model.Task{}, ErrInvalidTaskTitle
	}
	return s.repo.Create(ctx, userID, title)
}
