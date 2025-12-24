package handler

import (
	"errors"
	"net/http"

	"example.com/atode/backend/internal/auth"
	"example.com/atode/backend/internal/httpx"
	"example.com/atode/backend/internal/model"
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

type listTasksResponse struct {
	Tasks []model.Task `json:"tasks"`
}

func (h *TaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.List(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		httpx.WriteError(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
	}
}

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, r, http.StatusUnauthorized, "unauthorized", "unauthorized")
		return
	}

	tasks, err := h.tasks.ListTasks(r.Context(), userID)
	if err != nil {
		httpx.WriteError(w, r, http.StatusInternalServerError, "internal_server_error", "internal server error")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, listTasksResponse{Tasks: tasks})
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, r, http.StatusUnauthorized, "unauthorized", "unauthorized")
		return
	}

	var req createTaskRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteError(w, r, http.StatusBadRequest, "invalid_json", "invalid json")
		return
	}

	task, err := h.tasks.CreateTask(r.Context(), userID, req.Title)
	if err != nil {
		if errors.Is(err, service.ErrInvalidTaskTitle) {
			httpx.WriteError(w, r, http.StatusBadRequest, "title_required", "title is required")
			return
		}
		httpx.WriteError(w, r, http.StatusInternalServerError, "internal_server_error", "internal server error")
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, task)
}
