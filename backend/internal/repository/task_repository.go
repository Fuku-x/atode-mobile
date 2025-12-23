package repository

import (
	"context"
	"database/sql"

	"example.com/atode/backend/internal/model"
	"github.com/google/uuid"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, userID uuid.UUID, title string) (model.Task, error) {
	row := r.db.QueryRowContext(
		ctx,
		`INSERT INTO tasks (user_id, title)
VALUES ($1, $2)
RETURNING id, user_id, title, is_done, due_at, scheduled_at, created_at, updated_at;`,
		userID,
		title,
	)

	var t model.Task
	var dueAt sql.NullTime
	var scheduledAt sql.NullTime
	if err := row.Scan(
		&t.ID,
		&t.UserID,
		&t.Title,
		&t.IsDone,
		&dueAt,
		&scheduledAt,
		&t.CreatedAt,
		&t.UpdatedAt,
	); err != nil {
		return model.Task{}, err
	}

	if dueAt.Valid {
		t.DueAt = &dueAt.Time
	}
	if scheduledAt.Valid {
		t.ScheduledAt = &scheduledAt.Time
	}

	return t, nil
}
