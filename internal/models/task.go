package models

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

func NewTask(id int, title, description string) *Task {
	return &Task{ID: id, Title: title, Description: description}
}

func (t *Task) TaskDone() {
	t.Completed = true
}
