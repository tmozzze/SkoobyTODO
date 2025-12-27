package handlers

import (
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
		log.Warn("failed to decode reauest body", "err", err)
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid reauest body"})

		return
	}

	// create task
	id, err := h.service.Create(r.Context(), task)
	if err != nil {
		// validate error
		if errors.Is(err, service.ErrInvalidTitle) {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		// unknown error
		h.log.Error("internal error", "err", err)
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	// Sucsess
	respondJSON(w, http.StatusCreated, map[string]int{"id": id})

}

func (h *Handler) getAllTasks(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) getTaskByID(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.getTaskByID"
	// add operation to log
	log := h.log.With("op", op)

	// get id from URL
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Warn("invalid id path parametr", "id", idStr)
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	// getting task
	task, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		// invalid id error
		if errors.Is(err, service.ErrInvalidID) {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		// not found error
		if errors.Is(err, service.ErrTaskNotFound) {
			respondJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}

		// unknown error
		log.Error("internal error", "err", err)
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}
	// OK
	respondJSON(w, http.StatusOK, task)

}

func (h *Handler) updateTask(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) deleteTask(w http.ResponseWriter, r *http.Request) {

}

func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}
