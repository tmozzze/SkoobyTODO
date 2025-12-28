package service

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"reflect"
	"testing"

	"github.com/tmozzze/SkoobyTODO/internal/models"
)

// MockRepository - mock for tests
type MockRepository struct {
	CreateFunc  func(ctx context.Context, task models.Task) (int, error)
	UpdateFunc  func(ctx context.Context, id int, task models.Task) (models.Task, error)
	DeleteFunc  func(ctx context.Context, id int) error
	GetByIDFunc func(ctx context.Context, id int) (models.Task, error)
	GetAllFunc  func(ctx context.Context) ([]models.Task, error)
}

// Create - calls mock CreateFunc
func (m *MockRepository) Create(ctx context.Context, task models.Task) (int, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, task)
	}
	return 0, nil
}

// Delete - calls mock DeleteFunc
func (m *MockRepository) Delete(ctx context.Context, id int) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

// Update - calls mock UpdateFunc
func (m *MockRepository) Update(ctx context.Context, id int, task models.Task) (models.Task, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, id, task)
	}
	return models.Task{}, nil
}

// GetByID - calls mock GetByIDFunc
func (m *MockRepository) GetByID(ctx context.Context, id int) (models.Task, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return models.Task{}, nil
}

// GetAll - calls mock GetAllFunc
func (m *MockRepository) GetAll(ctx context.Context) ([]models.Task, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx)
	}
	return nil, nil
}

// - - - - TESTS - - - -

// TestTaskService_Create - test Create
func TestTaskService_Create(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	tests := []struct {
		name         string
		inputTask    models.Task
		mockBehavior func(m *MockRepository)
		wantErr      bool
		expectedID   int
	}{
		{
			name:      "Success",
			inputTask: models.Task{Title: "Test task"},
			mockBehavior: func(m *MockRepository) {
				m.CreateFunc = func(ctx context.Context, task models.Task) (int, error) {
					return 1, nil
				}
			},
			wantErr:    false,
			expectedID: 1,
		},
		{
			name:      "Validation Error - Empty Title",
			inputTask: models.Task{Title: ""}, // empty title
			mockBehavior: func(m *MockRepository) {
				// do nothing
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create mock
			mockRepo := &MockRepository{}
			// mock falling
			if tt.mockBehavior != nil {
				tt.mockBehavior(mockRepo)
			}

			// create service
			svc := NewService(mockRepo, log)

			// method
			gotID, err := svc.Create(context.Background(), tt.inputTask)

			// check results
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotID != tt.expectedID {
				t.Errorf("Create() gotID = %v, want %v", gotID, tt.expectedID)
			}
		})
	}

}

// TestTaskService_Update - test Update
func TestTaskService_Update(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	tests := []struct {
		name         string
		inputID      int
		inputTask    models.Task
		mockBehavior func(m *MockRepository)
		wantErr      bool
		checkErr     error // expected err
	}{
		{
			name:      "Success",
			inputID:   1,
			inputTask: models.Task{Title: "Updated Task"},
			mockBehavior: func(m *MockRepository) {
				m.UpdateFunc = func(ctx context.Context, id int, task models.Task) (models.Task, error) {
					task.ID = id
					return task, nil
				}
			},
			wantErr: false,
		},
		{
			name:      "Validation Error - Empty Title",
			inputID:   1,
			inputTask: models.Task{Title: ""},
			mockBehavior: func(m *MockRepository) {
				// do nothing
			},
			wantErr: true,
		},
		{
			name:      "Validation Error - Invalid ID",
			inputID:   0,
			inputTask: models.Task{Title: "Valid"},
			mockBehavior: func(m *MockRepository) {
				// do nothing
			},
			wantErr: true,
		},
		{
			name:      "Not Found",
			inputID:   999,
			inputTask: models.Task{Title: "Valid"},
			mockBehavior: func(m *MockRepository) {
				m.UpdateFunc = func(ctx context.Context, id int, task models.Task) (models.Task, error) {
					return models.Task{}, models.ErrTaskNotFound
				}
			},
			wantErr:  true,
			checkErr: models.ErrTaskNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRepository{}
			if tt.mockBehavior != nil {
				tt.mockBehavior(mockRepo)
			}

			svc := NewService(mockRepo, log)

			_, err := svc.Update(context.Background(), tt.inputID, tt.inputTask)

			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}

			// expected err
			if tt.checkErr != nil {
				if !errors.Is(err, tt.checkErr) {
					t.Errorf("Update() error = %v, want specific error %v", err, tt.checkErr)
				}
			}
		})
	}
}

// TestTaskService_Delete - test Delete
func TestTaskService_Delete(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	tests := []struct {
		name         string
		inputID      int
		mockBehavior func(m *MockRepository)
		wantErr      bool
	}{
		{
			name:    "Success",
			inputID: 1,
			mockBehavior: func(m *MockRepository) {
				m.DeleteFunc = func(ctx context.Context, id int) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name:    "Invalid ID",
			inputID: -1,
			mockBehavior: func(m *MockRepository) {
				// do nothing
			},
			wantErr: true,
		},
		{
			name:    "Not Found",
			inputID: 999,
			mockBehavior: func(m *MockRepository) {
				m.DeleteFunc = func(ctx context.Context, id int) error {
					return models.ErrTaskNotFound
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRepository{}
			if tt.mockBehavior != nil {
				tt.mockBehavior(mockRepo)
			}

			svc := NewService(mockRepo, log)

			err := svc.Delete(context.Background(), tt.inputID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestTaskService_GetByID - test GetByID
func TestTaskService_GetByID(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	tests := []struct {
		name         string
		inputID      int
		mockBehavior func(m *MockRepository)
		wantErr      bool
	}{
		{
			name:    "Success",
			inputID: 1,
			mockBehavior: func(m *MockRepository) {
				m.GetByIDFunc = func(ctx context.Context, id int) (models.Task, error) {
					return models.Task{ID: 1, Title: "Found"}, nil
				}
			},
			wantErr: false,
		},
		{
			name:    "Not Found",
			inputID: 2,
			mockBehavior: func(m *MockRepository) {
				m.GetByIDFunc = func(ctx context.Context, id int) (models.Task, error) {
					return models.Task{}, models.ErrTaskNotFound
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRepository{}
			if tt.mockBehavior != nil {
				tt.mockBehavior(mockRepo)
			}

			svc := NewService(mockRepo, log)

			_, err := svc.GetByID(context.Background(), tt.inputID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestTaskService_GetAll - test GetAll
func TestTaskService_GetAll(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	tests := []struct {
		name         string
		mockBehavior func(m *MockRepository)
		wantTasks    []models.Task
		wantErr      bool
	}{
		{
			name: "Success - Tasks Found",
			mockBehavior: func(m *MockRepository) {
				m.GetAllFunc = func(ctx context.Context) ([]models.Task, error) {
					return []models.Task{
						{ID: 1, Title: "Task 1", Completed: false},
						{ID: 2, Title: "Task 2", Completed: true},
					}, nil
				}
			},
			wantTasks: []models.Task{
				{ID: 1, Title: "Task 1", Completed: false},
				{ID: 2, Title: "Task 2", Completed: true},
			},
			wantErr: false,
		},
		{
			name: "Success - Empty Slice",
			mockBehavior: func(m *MockRepository) {
				m.GetAllFunc = func(ctx context.Context) ([]models.Task, error) {
					// empty slice
					return []models.Task{}, nil
				}
			},
			wantTasks: []models.Task{}, // empty slice
			wantErr:   false,
		},
		{
			name: "Repository Error",
			mockBehavior: func(m *MockRepository) {
				m.GetAllFunc = func(ctx context.Context) ([]models.Task, error) {
					return nil, errors.New("db failure")
				}
			},
			wantTasks: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRepository{}
			if tt.mockBehavior != nil {
				tt.mockBehavior(mockRepo)
			}

			svc := NewService(mockRepo, log)

			gotTasks, err := svc.GetAll(context.Background())

			// 1. Check err
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 2. Check slices
			if !reflect.DeepEqual(gotTasks, tt.wantTasks) {
				t.Errorf("GetAll() got = %v, want %v", gotTasks, tt.wantTasks)
			}
		})
	}
}
