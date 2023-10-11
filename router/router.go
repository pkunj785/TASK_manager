package router

import (
	"com.serve_volt/middleware"
	"github.com/go-chi/chi/v5"
)

func Routers(r chi.Router) {
	r.Get("/api/getTasks", middleware.GetAllTasks)
	r.Post("/api/createTask", middleware.CreateTask)
	r.Put("/api/task/{id}", middleware.TaskComplete)
	r.Put("/api/undoTask/{id}", middleware.UndoTask)
	r.Delete("/api/deleteTask/{id}", middleware.DeleteTask)
	r.Delete("/api/deleteAllTasks", middleware.DeleteAllTasks)
}
