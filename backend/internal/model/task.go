package model

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"userId"`
	Title       string     `json:"title"`
	IsDone      bool       `json:"isDone"`
	DueAt       *time.Time `json:"dueAt,omitempty"`
	ScheduledAt *time.Time `json:"scheduledAt,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}
