package inmemory

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/tmozzze/SkoobyTODO/internal/models"
)

// TestMemStorage_Create - test storage Create
func TestMemStorage_Create(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	storage := NewMemStorage(log)
	ctx := context.Background()

	// TASK 1
	task1 := models.Task{
		Title:       "Test Task 1 joj",
		Description: "This is a test task",
	}
	id1, err := storage.Create(ctx, task1)
	if err != nil {
		t.Fatalf("Failed to create task 1: %v", err)
	}
	if id1 != 1 {
		t.Fatalf("Expected task ID 1, got %d", id1)
	}

	// TASK 2
	task2 := models.Task{
		Title:       "Test Task 2 sas",
		Description: "Another test task",
	}
	id2, err := storage.Create(ctx, task2)
	if err != nil {
		t.Fatalf("Failed to create task 2: %v", err)
	}
	if id2 != 2 {
		t.Fatalf("Expected task id 2 got %d", id2)
	}

	savedTask, err := storage.GetByID(ctx, id1)
	if err != nil {
		t.Fatalf("Failed to get task")
	}
	if savedTask.Title != task1.Title || savedTask.Description != task1.Description {
		t.Errorf("Saved title or description mismatch")
	}
}

// TestMemStorage_Create - test storage duplicate Create
func TestMemStorage_Create_Duplication(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	storage := NewMemStorage(log)
	ctx := context.Background()

	// take ID 1 and return the counter back to 0
	storage.store[1] = models.Task{ID: 1, Title: "I was here first"}
	storage.lastID = 0

	// call Create
	// nextID = lastID + 1 (0 + 1 = 1)
	// storage should see that its taken
	_, err := storage.Create(ctx, models.Task{Title: "New Task"})
	if err == nil {
		t.Fatal("Expected error due to ID duplication, got nil")
	}
}

// TestMemStorage_GetAll_Sorted - test storage GetAll
func TestMemStorage_GetAll_Sorted(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	storage := NewMemStorage(log)
	ctx := context.Background()

	storage.Create(ctx, models.Task{Title: "Sas 1"})
	storage.Create(ctx, models.Task{Title: "Sas 2"})

	tasks, _ := storage.GetAll(ctx)

	if len(tasks) != 2 {
		t.Fatalf("Expected 2 tasks")
	}

	if tasks[0].ID > tasks[1].ID {
		t.Errorf("Task are not sorted")
	}
}

// TestMemStorage_Update - test storage Update
func TestMemStorage_Update(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	storage := NewMemStorage(log)
	ctx := context.Background()

	task := models.Task{Title: "sas", Description: "Test description"}
	updTask := models.Task{Title: "sas", Description: "Test upd description"}

	id, err := storage.Create(ctx, task)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	_, err = storage.Update(ctx, id, updTask)
	if err != nil {
		t.Fatalf("Failed to update task: %v", err)
	}

	resTask, err := storage.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("Failed to get task")
	}
	if resTask.Title != updTask.Title || resTask.Description != updTask.Description {
		t.Errorf("Title or description mismatch")
	}

}
