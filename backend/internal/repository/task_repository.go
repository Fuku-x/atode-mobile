package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"example.com/atode/backend/internal/model"
	"github.com/google/uuid"
)

type TaskPatch struct {
	Title             *string
	IsDone            *bool
	DueAtPresent      bool
	DueAt             *time.Time
	ScheduledAtPresent bool
	ScheduledAt        *time.Time
}

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

func (r *TaskRepository) Delete(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) (bool, error) {
	res, err := r.db.ExecContext(
		ctx,
		`DELETE FROM tasks WHERE id = $1 AND user_id = $2;`,
		taskID,
		userID,
	)
	if err != nil {
		return false, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return n > 0, nil
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

func (r *TaskRepository) Update(ctx context.Context, userID uuid.UUID, taskID uuid.UUID, patch TaskPatch) (model.Task, error) {
	setClauses := make([]string, 0, 5)
	args := make([]any, 0, 7)
	idx := 1

	if patch.Title != nil {
		setClauses = append(setClauses, fmt.Sprintf("title = $%d", idx))
		args = append(args, *patch.Title)
		idx++
	}
	if patch.IsDone != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_done = $%d", idx))
		args = append(args, *patch.IsDone)
		idx++
	}
	if patch.DueAtPresent {
		setClauses = append(setClauses, fmt.Sprintf("due_at = $%d", idx))
		if patch.DueAt != nil {
			args = append(args, *patch.DueAt)
		} else {
			args = append(args, nil)
		}
		idx++
	}
	if patch.ScheduledAtPresent {
		setClauses = append(setClauses, fmt.Sprintf("scheduled_at = $%d", idx))
		if patch.ScheduledAt != nil {
			args = append(args, *patch.ScheduledAt)
		} else {
			args = append(args, nil)
		}
		idx++
	}

	setClauses = append(setClauses, "updated_at = now()")

	args = append(args, taskID, userID)

	query := `UPDATE tasks
SET ` + strings.Join(setClauses, ", ") + fmt.Sprintf("\nWHERE id = $%d AND user_id = $%d\n", idx, idx+1) +
		`RETURNING id, user_id, title, is_done, due_at, scheduled_at, created_at, updated_at;`

	row := r.db.QueryRowContext(ctx, query, args...)

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
