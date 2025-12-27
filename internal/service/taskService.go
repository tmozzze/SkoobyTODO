package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/tmozzze/SkoobyTODO/internal/models"
)

type TaskRepository interface {
	Create(ctx context.Context, task models.Task) (id int, err error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, id int, updTask models.Task) (task models.Task, err error)
	GetByID(ctx context.Context, id int) (task models.Task, err error)
	GetAll(ctx context.Context) (tasks []models.Task, err error)
}

// TaskService - struct for service
type TaskService struct {
	repo TaskRepository
	log  *slog.Logger
}

// NewService - constructor for TaskService --> *TaskService
func NewService(taskRepo TaskRepository, log *slog.Logger) *TaskService {
	return &TaskService{repo: taskRepo, log: log}
}

// Create - creates task --> id, error
func (s *TaskService) Create(ctx context.Context, task models.Task) (int, error) {
	const op = "service.Create"
	// add operation to log
	log := s.log.With("op", op)

	log.Debug("starting to create task")

	// Validating
	if task.Title == "" {
		log.Warn("empty title")

		return 0, fmt.Errorf("%s: %w", op, ErrInvalidTitle)
	}

	id, err := s.repo.Create(ctx, task)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("task created", "id", id)

	return id, nil

}

// Delete - deletes task by id --> error
func (s *TaskService) Delete(ctx context.Context, id int) error {
	const op = "service.Delete"
	// add operation to log
	log := s.log.With("op", op)

	log.Debug("starting to delete task")

	// Validating
	if id <= 0 {
		log.Warn("invalid id")
		return fmt.Errorf("%s: %w", op, ErrInvalidID)
	}

	err := s.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("task deleted", "id", id)

	return nil
}

// Update - update task to updTask by id --> updTask, error
func (s *TaskService) Update(ctx context.Context, id int, updTask models.Task) (models.Task, error) {
	const op = "service.Update"
	// add operation to log
	log := s.log.With("op", op)

	log.Debug("starting to update task")

	// Validating
	if id <= 0 {
		log.Warn("invalid id")
		return models.Task{}, fmt.Errorf("%s: %w", op, ErrInvalidID)
	}

	if updTask.Title == "" {
		log.Warn("empty title")
		return models.Task{}, fmt.Errorf("%s: %w", op, ErrInvalidTitle)
	}

	task, err := s.repo.Update(ctx, id, updTask)
	if err != nil {
		return models.Task{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("task updated", "id", id)

	return task, nil

}

// GetById - get task by id --> task, error
func (s *TaskService) GetByID(ctx context.Context, id int) (models.Task, error) {
	const op = "service.GetByID"
	// add operation to log
	log := s.log.With("op", op)

	log.Debug("starting to get task by id")

	// Validating
	if id <= 0 {
		log.Warn("invalid id")
		return models.Task{}, fmt.Errorf("%s: %w", op, ErrInvalidID)
	}

	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return models.Task{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("got task")

	return task, nil
}

// GetAll - get all tasks --> []task, error
func (s *TaskService) GetAll(ctx context.Context) ([]models.Task, error) {
	const op = "service.GetAll"
	// add operation to log
	log := s.log.With("op", op)

	log.Debug("starting to get all tasks")

	tasks, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("get tasks", "amount", len(tasks))

	return tasks, nil

}
