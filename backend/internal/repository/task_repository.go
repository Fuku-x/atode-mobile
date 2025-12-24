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

func (r *TaskRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]model.Task, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, title, is_done, due_at, scheduled_at, created_at, updated_at
FROM tasks
WHERE user_id = $1
ORDER BY is_done ASC, COALESCE(due_at, 'infinity'::timestamptz) ASC, created_at DESC;`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.Task
	for rows.Next() {
		var t model.Task
		var dueAt sql.NullTime
		var scheduledAt sql.NullTime
		if err := rows.Scan(
			&t.ID,
			&t.UserID,
			&t.Title,
			&t.IsDone,
			&dueAt,
			&scheduledAt,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if dueAt.Valid {
			t.DueAt = &dueAt.Time
		}
		if scheduledAt.Valid {
			t.ScheduledAt = &scheduledAt.Time
		}
		out = append(out, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}
