package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Wimer/TicketControlSystem/internal/repository"
	"github.com/Wimer/TicketControlSystem/internal/service"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

// TaskHandler - обработчик HTTP-запросов
type TaskHandler struct {
	service *service.TaskService
	logger  zerolog.Logger
	apiKey  string
}

// NewTaskHandler - конструктор обработчика
func NewTaskHandler(service *service.TaskService, logger zerolog.Logger, apiKey string) *TaskHandler {
	return &TaskHandler{service: service, logger: logger, apiKey: apiKey}
}

// AuthMiddleware - проверка API-ключа
func (h *TaskHandler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey != h.apiKey {
			h.logger.Warn().Msg("Unauthorized access attempt")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

// CreateTask - создание задачи
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	h.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Title  string `json:"title"`
			Status string `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			h.logger.Error().Err(err).Msg("Failed to decode request")
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		task, err := h.service.CreateTask(input.Title, input.Status)
		if err != nil {
			h.logger.Error().Err(err).Msg("Failed to create task")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)
		h.logger.Info().Msg("Task created successfully")
	})(w, r)
}

// GetTask - получение задачи
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	h.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			h.logger.Error().Err(err).Msg("Invalid task ID")
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		task, err := h.service.GetTask(id)
		if err != nil {
			h.logger.Error().Err(err).Msg("Failed to get task")
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)
		h.logger.Info().Msg("Task retrieved successfully")
	})(w, r)
}

// UpdateTask - обновление задачи
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	h.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			h.logger.Error().Err(err).Msg("Invalid task ID")
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		var input struct {
			Title  string `json:"title"`
			Status string `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			h.logger.Error().Err(err).Msg("Failed to decode request")
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		if err := h.service.UpdateTask(id, input.Title, input.Status); err != nil {
			h.logger.Error().Err(err).Msg("Failed to update task")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		h.logger.Info().Msg("Task updated successfully")
	})(w, r)
}

// DeleteTask - удаление задачи
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	h.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			h.logger.Error().Err(err).Msg("Invalid task ID")
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		if err := h.service.DeleteTask(id); err != nil {
			h.logger.Error().Err(err).Msg("Failed to delete task")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		h.logger.Info().Msg("Task deleted successfully")
	})(w, r)
}
