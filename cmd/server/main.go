package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/Engineer-DF/task-s/internal/api"
	"github.com/Engineer-DF/task-s/internal/domain"
	"github.com/Engineer-DF/task-s/internal/repository"
	"github.com/Engineer-DF/task-s/internal/service"
)

type Server struct {
	log    *slog.Logger
	router *http.ServeMux
}

func NewServer(l *slog.Logger, taskService *service.Service) *Server {
	mux := http.NewServeMux()

	handler := api.NewTaskHandler(taskService, l)

	mux.HandleFunc("GET /tasks", handler.GetAll)
	mux.HandleFunc("GET /tasks/{id}", handler.GetByID)
	mux.HandleFunc("POST /tasks", handler.Create)
	mux.HandleFunc("PUT /tasks/{id}", handler.Update)
	mux.HandleFunc("DELETE /tasks/{id}", handler.Delete)

	s := &Server{
		log:    l,
		router: mux,
	}

	return s
}

func main() {
	l := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	initialTasks := make(map[int64]domain.Task)
	repo := repository.NewTaskRepository(initialTasks)
	taskService := service.NewService(repo)

	s := NewServer(l, taskService)

	s.log.Info("HTTP server started", "port", 8080)

	if err := http.ListenAndServe(":8080", s.router); err != nil {
		s.log.Error("HTTP server stopped", "error", err)
	}
}
