package inmemory

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"sync"

	"github.com/tmozzze/SkoobyTODO/internal/models"
)

// MemStorage - struct to async caching tasks
type MemStorage struct {
	mu     sync.RWMutex
	store  map[int]models.Task
	lastID int
	log    *slog.Logger
}

// NewMemStorage - initialize inmemory storage
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
	// Add operation to log
	log := m.log.With(slog.String("op", op))

	log.Debug("starting to create task")

	select {
	case <-ctx.Done():
		log.Debug("context done")
		return 0, ctx.Err()
	default:

	}

	m.mu.Lock()
	defer m.mu.Unlock()

	id := m.lastID + 1
	task.ID = id

	if _, exist := m.store[id]; exist {
		log.Warn("task already exist", "id", id)
		return 0, fmt.Errorf("%s: task id already exist", op)
	}
	m.store[id] = task
	m.lastID = id

	log.Debug("task created", "task", task)

	return id, nil
}

// Delete - deletes task from inmemory storage by id --> error
func (m *MemStorage) Delete(ctx context.Context, id int) error {
	const op = "storage.inmemory.Delete"
	// Add operation to log
	log := m.log.With(slog.String("op", op))

	log.Debug("starting to delete task", "id", id)

	select {
	case <-ctx.Done():
		log.Debug("context done")
		return ctx.Err()
	default:
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	task, exists := m.store[id]
	if !exists {
		log.Warn("task doesn't exist", "id", id)
		return fmt.Errorf("%s: %w", op, ErrTaskNotFound)
	}

	delete(m.store, id)

	log.Debug("task deleted", "task", task)

	return nil
}

// GetByID - get task by id from inmemory storage --> task, error
func (m *MemStorage) GetByID(ctx context.Context, id int) (models.Task, error) {
	const op = "storage.inmemory.GetByID"
	// Add operation to log
	log := m.log.With(slog.String("op", op))

	log.Debug("starting to get task", "id", id)

	select {
	case <-ctx.Done():
		log.Debug("context done")
		return models.Task{}, ctx.Err()
	default:
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	task, exists := m.store[id]
	if !exists {
		log.Warn("task doesn't exist", "id", id)
		return models.Task{}, fmt.Errorf("%s: %w", op, ErrTaskNotFound)
	}

	log.Debug("got task", "task", task)

	return task, nil

}

// Update - updates the task in inmemory storage --> task, error
func (m *MemStorage) Update(ctx context.Context, id int, updTask models.Task) (models.Task, error) {
	const op = "storage.inmemory.Update"
	// Add operation to log
	log := m.log.With(slog.String("op", op))

	log.Debug("starting to update task", "id", id)

	select {
	case <-ctx.Done():
		log.Debug("context done")
		return models.Task{}, ctx.Err()
	default:
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	_, exists := m.store[id]
	if !exists {
		log.Warn("task doesn't exist")
		return models.Task{}, fmt.Errorf("%s: %w", op, ErrTaskNotFound)
	}

	updTask.ID = id

	m.store[id] = updTask

	log.Debug("task updated", "task", updTask)

	return updTask, nil

}

// GetAll - get all task from inmemory storage --> []task, error
func (m *MemStorage) GetAll(ctx context.Context) ([]models.Task, error) {
	const op = "storage.inmemory.GetAll"
	// Add operation to log
	log := m.log.With(slog.String("op", op))

	log.Debug("starting to get all task")

	select {
	case <-ctx.Done():
		log.Debug("context done")
		return nil, ctx.Err()
	default:
	}

	var tasks []models.Task

	func() {
		m.mu.RLock()
		defer m.mu.RUnlock()

		tasks = make([]models.Task, 0, len(m.store))
		for _, task := range m.store {
			tasks = append(tasks, task)
		}

	}()

	select {
	case <-ctx.Done():
		log.Debug("context done")
		return nil, ctx.Err()
	default:
	}

	sort.Slice(tasks, func(i, j int) bool { return tasks[i].ID < tasks[j].ID })

	if len(tasks) <= 0 {
		log.Warn("list of tasks is empty")
	}

	log.Debug("got all tasks", "tasks", tasks)

	return tasks, nil

}
