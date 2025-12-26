package inmemory

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/tmozzze/SkoobyTODO/internal/models"
)

type MemStorage struct {
	mu     sync.RWMutex
	store  map[int]models.Task
	lastID int
	log    *slog.Logger
}

func NewMemStorage(log *slog.Logger) *MemStorage {
	return &MemStorage{
		store:  make(map[int]models.Task),
		lastID: 0,
		log:    log,
	}
}

// Create - creates task in inmemory storage --> id, error
func (m *MemStorage) Create(ctx context.Context, task models.Task) (int, error) {
	const op = "storage.inmemory.Create"

	memLog := m.log.With(slog.String("op", op))
	memLog.Debug("starting to create task")

	select {
	case <-ctx.Done():
		memLog.Debug("context done")
		return 0, ctx.Err()
	default:

	}

	m.mu.Lock()
	defer m.mu.Unlock()

	id := m.lastID + 1
	task.ID = id

	if _, exist := m.store[id]; exist {
		memLog.Warn("task already exist", "id", id)
		return 0, fmt.Errorf("%s: task id already exist", op)
	}
	m.store[id] = task
	m.lastID = id

	memLog.Info("task created", "id", id)
	memLog.Debug("task created", "task", task)

	return id, nil
}

// Delete - deletes task from inmemory storage by id --> error
func (m *MemStorage) Delete(ctx context.Context, id int) error {
	const op = "storage.inmemory.Delete"

	memLog := m.log.With(slog.String("op", op))
	memLog.Debug("starting to delete task", "id", id)

	select {
	case <-ctx.Done():
		memLog.Debug("context done")
		return ctx.Err()
	default:
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	task, exists := m.store[id]
	if !exists {
		memLog.Warn("task doesn't exist", "id", id)
		return fmt.Errorf("%s: %w", op, ErrTaskNotFound)
	}

	delete(m.store, id)
	memLog.Info("task deleted", "id", id)
	memLog.Debug("task deleted", "task", task)

	return nil
}

// GetByID - get task by id --> task, error
func (m *MemStorage) GetByID(ctx context.Context, id int) (models.Task, error) {
	const op = "storage.inmemory.GetByID"

	memLog := m.log.With(slog.String("op", op))
	memLog.Debug("getting task", "id", id)

	select {
	case <-ctx.Done():
		memLog.Debug("context done")
		return models.Task{}, ctx.Err()
	default:
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	task, exists := m.store[id]
	if !exists {
		memLog.Warn("task not exists", "id", id)
		return models.Task{}, fmt.Errorf("%s: %w", op, ErrTaskNotFound)
	}

	memLog.Info("got task", "id", id)
	memLog.Debug("got task", "task", task)

	return task, nil

}

// Update - updates
