package handlers

import "net/http"

func (h *Handler) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /todos", h.createTask)
	mux.HandleFunc("GET /todos", h.getAllTasks)
	mux.HandleFunc("GET /todos/{id}", h.getTaskByID)
	mux.HandleFunc("PUT /todos/{id}", h.updateTask)
	mux.HandleFunc("DELETE /todos/{id}", h.deleteTask)

	return mux
}
