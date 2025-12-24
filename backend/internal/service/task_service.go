package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"example.com/atode/backend/internal/model"
	"example.com/atode/backend/internal/repository"
	"github.com/google/uuid"
)

var ErrInvalidTaskTitle = errors.New("invalid task title")
var ErrNoTaskFieldsToUpdate = errors.New("no task fields to update")

type TaskCreator interface {
	Create(ctx context.Context, userID uuid.UUID, title string) (model.Task, error)
}

type TaskLister interface {
	ListByUser(ctx context.Context, userID uuid.UUID) ([]model.Task, error)
}

type TaskUpdater interface {
	Update(ctx context.Context, userID uuid.UUID, taskID uuid.UUID, patch repository.TaskPatch) (model.Task, error)
}

type TaskService struct {
	repo interface {
		TaskCreator
		TaskLister
		TaskUpdater
	}
}

func NewTaskService(repo interface {
	TaskCreator
	TaskLister
	TaskUpdater
}) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(ctx context.Context, userID uuid.UUID, title string) (model.Task, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return model.Task{}, ErrInvalidTaskTitle
	}
	return s.repo.Create(ctx, userID, title)
}

func (s *TaskService) ListTasks(ctx context.Context, userID uuid.UUID) ([]model.Task, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *TaskService) UpdateTask(
	ctx context.Context,
	userID uuid.UUID,
	taskID uuid.UUID,
	title *string,
	isDone *bool,
	dueAtPresent bool,
	dueAt *time.Time,
	scheduledAtPresent bool,
	scheduledAt *time.Time,
) (model.Task, error) {
	if title != nil {
		trimmed := strings.TrimSpace(*title)
		if trimmed == "" {
			return model.Task{}, ErrInvalidTaskTitle
		}
		title = &trimmed
	}

	patch := repository.TaskPatch{
		Title:              title,
		IsDone:             isDone,
		DueAtPresent:       dueAtPresent,
		DueAt:              dueAt,
		ScheduledAtPresent: scheduledAtPresent,
		ScheduledAt:        scheduledAt,
	}

	if patch.Title == nil && patch.IsDone == nil && !patch.DueAtPresent && !patch.ScheduledAtPresent {
		return model.Task{}, ErrNoTaskFieldsToUpdate
	}

	return s.repo.Update(ctx, userID, taskID, patch)
}
