package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"example.com/atode/backend/internal/auth"
	"example.com/atode/backend/internal/httpx"
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
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, r, http.StatusUnauthorized, "unauthorized", "unauthorized")
		return
	}

	taskID, status, ok := parseTaskIDFromPath(r.URL.Path)
	if !ok {
		if status == http.StatusBadRequest {
			httpx.WriteError(w, r, http.StatusBadRequest, "invalid_task_id", "invalid task id")
			return
		}
		httpx.WriteError(w, r, http.StatusNotFound, "not_found", "not found")
		return
	}

	switch r.Method {
	case http.MethodPut:
		h.handlePut(w, r, userID, taskID)
	case http.MethodDelete:
		h.handleDelete(w, r, userID, taskID)
	default:
		httpx.WriteError(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
	}
}


func parseTaskIDFromPath(path string) (uuid.UUID, int, bool) {
	idStr := strings.TrimPrefix(path, "/tasks/")
	if idStr == "" || strings.Contains(idStr, "/") {
		return uuid.UUID{}, http.StatusNotFound, false
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.UUID{}, http.StatusBadRequest, false
	}
	return id, http.StatusOK, true
}

func (h *TaskItemHandler) handleDelete(w http.ResponseWriter, r *http.Request, userID uuid.UUID, taskID uuid.UUID) {
	if err := h.tasks.DeleteTask(r.Context(), userID, taskID); err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			httpx.WriteError(w, r, http.StatusNotFound, "not_found", "not found")
			return
		}
		httpx.WriteError(w, r, http.StatusInternalServerError, "internal_server_error", "internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskItemHandler) handlePut(w http.ResponseWriter, r *http.Request, userID uuid.UUID, taskID uuid.UUID) {
	// Decode as raw map first to detect presence (including explicit null) and reject unknown fields.
	var raw map[string]json.RawMessage
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&raw); err != nil {
		httpx.WriteError(w, r, http.StatusBadRequest, "invalid_json", "invalid json")
		return
	}
	for k := range raw {
		switch k {
		case "title", "isDone", "dueAt", "scheduledAt":
		default:
			httpx.WriteError(w, r, http.StatusBadRequest, "invalid_json", "invalid json")
			return
		}
	}

	var (
		title              *string
		isDone             *bool
		dueAtPresent       bool
		dueAt              *time.Time
		scheduledAtPresent bool
		scheduledAt        *time.Time
	)

	if b, ok := raw["title"]; ok {
		var v string
		if err := json.Unmarshal(b, &v); err != nil {
			httpx.WriteError(w, r, http.StatusBadRequest, "invalid_title", "invalid title")
			return
		}
		title = &v
	}

	if b, ok := raw["isDone"]; ok {
		var v bool
		if err := json.Unmarshal(b, &v); err != nil {
			httpx.WriteError(w, r, http.StatusBadRequest, "invalid_is_done", "invalid isDone")
			return
		}
		isDone = &v
	}

	if b, ok := raw["dueAt"]; ok {
		dueAtPresent = true
		if string(b) != "null" {
			var s string
			if err := json.Unmarshal(b, &s); err != nil {
				httpx.WriteError(w, r, http.StatusBadRequest, "invalid_due_at", "invalid dueAt")
				return
			}
			t, err := time.Parse(time.RFC3339Nano, s)
			if err != nil {
				httpx.WriteError(w, r, http.StatusBadRequest, "invalid_due_at", "invalid dueAt")
				return
			}
			dueAt = &t
		}
	}

	if b, ok := raw["scheduledAt"]; ok {
		scheduledAtPresent = true
		if string(b) != "null" {
			var s string
			if err := json.Unmarshal(b, &s); err != nil {
				httpx.WriteError(w, r, http.StatusBadRequest, "invalid_scheduled_at", "invalid scheduledAt")
				return
			}
			t, err := time.Parse(time.RFC3339Nano, s)
			if err != nil {
				httpx.WriteError(w, r, http.StatusBadRequest, "invalid_scheduled_at", "invalid scheduledAt")
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
			httpx.WriteError(w, r, http.StatusBadRequest, "title_required", "title is required")
			return
		}
		if errors.Is(err, service.ErrNoTaskFieldsToUpdate) {
			httpx.WriteError(w, r, http.StatusBadRequest, "no_fields", "no fields to update")
			return
		}
		if errors.Is(err, sql.ErrNoRows) {
			httpx.WriteError(w, r, http.StatusNotFound, "not_found", "not found")
			return
		}

		httpx.WriteError(w, r, http.StatusInternalServerError, "internal_server_error", "internal server error")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, task)
}
