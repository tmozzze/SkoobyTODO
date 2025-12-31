package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/tmozzze/SkoobyTODO/internal/models"
	"github.com/tmozzze/SkoobyTODO/internal/service"
)

// Handler - struct for Handler
type Handler struct {
	service *service.TaskService
	log     *slog.Logger
}

// NewHandler - constructor for Handler --> *Handler
func NewHandler(service *service.TaskService, log *slog.Logger) *Handler {
	return &Handler{service: service, log: log}
}

func (h *Handler) createTask(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.createTask"
	// add operation to log
	log := h.log.With("op", op)

	var task models.Task

	// decode json to models.task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		// 400 Bad Request
		log.Warn("failed to decode request body", "err", err)
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})

		return
	}

	// create task (go to service)
	id, err := h.service.Create(r.Context(), task)
	if err != nil {
		h.handleError(w, err, op)
		return
	}

	// 201 Created
	respondJSON(w, http.StatusCreated, map[string]int{"id": id})

}

func (h *Handler) getAllTasks(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.getAllTasks"

	// get tasks (go to service)
	tasks, err := h.service.GetAll(r.Context())
	if err != nil {
		h.handleError(w, err, op)
		return
	}

	// 200 OK
	respondJSON(w, http.StatusOK, tasks)
}

func (h *Handler) getTaskByID(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.getTaskByID"
	// add operation to log
	log := h.log.With("op", op)

	// get id from URL
	idStr := r.PathValue("id")

	// convert id to int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// 400 Bad Request
		log.Warn("invalid id path parametr", "id", idStr)
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	// getting task (go to service)
	task, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.handleError(w, err, op)
		return
	}
	// 200 OK
	respondJSON(w, http.StatusOK, task)

}

func (h *Handler) updateTask(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.updateTask"
	// add operation to log
	log := h.log.With("op", op)

	// get id from URL
	idStr := r.PathValue("id")

	// convert id to int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// 400 Bad Request
		log.Warn("invalid id path parametr", "id", idStr)
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	// decode json to models.task
	var updTask models.Task

	err = json.NewDecoder(r.Body).Decode(&updTask)
	if err != nil {
		// 400 Bad Request
		log.Warn("failed to decode reques body", "err", err)
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	// update task (go to service)
	task, err := h.service.Update(r.Context(), id, updTask)
	if err != nil {
		h.handleError(w, err, op)
		return
	}

	// 200 OK
	respondJSON(w, http.StatusOK, task)
}

func (h *Handler) deleteTask(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.deleteTask"
	// add operation to log
	log := h.log.With("op", op)

	// get id from URL
	idStr := r.PathValue("id")

	// convert id to int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// 400 Bad Request
		log.Warn("invalid id path parametr", "id", idStr)
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	// delete task (go to service)
	err = h.service.Delete(r.Context(), id)
	if err != nil {
		h.handleError(w, err, op)
		return
	}

	// 204 No Content
	respondJSON(w, http.StatusNoContent, nil)

}

func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}

func (h *Handler) handleError(w http.ResponseWriter, err error, op string) {
	log := h.log.With("op", op)

	// 400 Bad Request (invalid id)
	if errors.Is(err, models.ErrInvalidID) {
		log.Warn("validation error", "err", err)
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": models.ErrInvalidID.Error()})
		return
	}
	// 400 Bad Request (invalid title)
	if errors.Is(err, models.ErrInvalidTitle) {
		log.Warn("validation error", "err", err)
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": models.ErrInvalidTitle.Error()})
		return
	}

	// 404 Not Found (task not found)
	if errors.Is(err, models.ErrTaskNotFound) {
		log.Warn("not found", "err", err)
		respondJSON(w, http.StatusNotFound, map[string]string{"error": models.ErrTaskNotFound.Error()})
		return
	}

	// timeout
	if errors.Is(err, context.DeadlineExceeded) {
		log.Warn("operation timed out", "err", err)
		return
	}

	// client canceled
	if errors.Is(err, context.Canceled) {
		log.Debug("client canceled request")
		return
	}

	// 500 Internal Error
	log.Error("internal error", "err", err)
	respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})

}
