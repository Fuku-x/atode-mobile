package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"example.com/atode/backend/internal/auth"
	"example.com/atode/backend/internal/service"
	"github.com/google/uuid"
)

type TaskItemHandler struct {
	tasks *service.TaskService
}

func NewTaskItemHandler(tasks *service.TaskService) *TaskItemHandler {
	return &TaskItemHandler{tasks: tasks}
}

type updateTaskRequest struct {
	Title       json.RawMessage `json:"title"`
	IsDone      json.RawMessage `json:"isDone"`
	DueAt       json.RawMessage `json:"dueAt"`
	ScheduledAt json.RawMessage `json:"scheduledAt"`
}

func (h *TaskItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	if idStr == "" || strings.Contains(idStr, "/") {
		writeJSON(w, http.StatusNotFound, errorResponse{Error: "not_found"})
		return
	}

	taskID, err := uuid.Parse(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid_task_id"})
		return
	}

	var req updateTaskRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid_json"})
		return
	}

	var (
		title              *string
		isDone             *bool
		dueAtPresent       bool
		dueAt              *time.Time
		scheduledAtPresent bool
		scheduledAt        *time.Time
	)

	if req.Title != nil {
		var v string
		if err := json.Unmarshal(req.Title, &v); err != nil {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid_title"})
			return
		}
		title = &v
	}

	if req.IsDone != nil {
		var v bool
		if err := json.Unmarshal(req.IsDone, &v); err != nil {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid_is_done"})
			return
		}
		isDone = &v
	}

	if req.DueAt != nil {
		dueAtPresent = true
		if string(req.DueAt) != "null" {
			var s string
			if err := json.Unmarshal(req.DueAt, &s); err != nil {
				writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid_due_at"})
				return
			}
			t, err := time.Parse(time.RFC3339Nano, s)
			if err != nil {
				writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid_due_at"})
				return
			}
			dueAt = &t
		}
	}

	if req.ScheduledAt != nil {
		scheduledAtPresent = true
		if string(req.ScheduledAt) != "null" {
			var s string
			if err := json.Unmarshal(req.ScheduledAt, &s); err != nil {
				writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid_scheduled_at"})
				return
			}
			t, err := time.Parse(time.RFC3339Nano, s)
			if err != nil {
				writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid_scheduled_at"})
				return
			}
			scheduledAt = &t
		}
	}

	task, err := h.tasks.UpdateTask(
		r.Context(),
		userID,
		taskID,
		title,
		isDone,
		dueAtPresent,
		dueAt,
		scheduledAtPresent,
		scheduledAt,
	)
	if err != nil {
		if errors.Is(err, service.ErrInvalidTaskTitle) {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "title_required"})
			return
		}
		if errors.Is(err, service.ErrNoTaskFieldsToUpdate) {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "no_fields"})
			return
		}
		if errors.Is(err, sql.ErrNoRows) {
			writeJSON(w, http.StatusNotFound, errorResponse{Error: "not_found"})
			return
		}

		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal_server_error"})
		return
	}

	writeJSON(w, http.StatusOK, task)
}
