package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"example.com/atode/backend/internal/auth"
	"example.com/atode/backend/internal/service"
)

type TaskHandler struct {
	tasks *service.TaskService
}

func NewTaskHandler(tasks *service.TaskService) *TaskHandler {
	return &TaskHandler{tasks: tasks}
}

type createTaskRequest struct {
	Title string `json:"title"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid_json"})
		return
	}

	task, err := h.tasks.CreateTask(r.Context(), userID, req.Title)
	if err != nil {
		if errors.Is(err, service.ErrInvalidTaskTitle) {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "title_required"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal_server_error"})
		return
	}

	writeJSON(w, http.StatusCreated, task)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
