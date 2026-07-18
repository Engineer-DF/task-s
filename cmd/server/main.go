package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
)

// TODO:
// Add multithreading support (goroutines, mutex, etc.)

//TODO:
// Refactor code for Clean Architectures!!!!!

type Server struct {
	tasks   map[int64]Task
	counter atomic.Int64
	log     *slog.Logger
}	

type Task struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Done        bool   `json:"done"`
}

func (s *Server) findMaxID() int64 {
	var maxID int64 = -1

	for key := range s.tasks {
		if key > maxID {
			maxID = key
		}
	}
	return maxID
}

func NewServer(m map[int64]Task, l *slog.Logger) *Server {

	s := &Server{
		tasks: m,
		log:   l,
	}
	s.counter.Store(s.findMaxID())
	return s
}

func (s *Server) ListTasks(w http.ResponseWriter, r *http.Request) {
	payload, err := json.Marshal(&s.tasks)
	if err != nil {
		s.log.Error("failed to marshal JSON", "error", err, "remote_addr", r.RemoteAddr)
		http.Error(w, "failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func (s *Server) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		s.log.Warn("invalid id format")
		http.Error(w, "invalid ID format", http.StatusBadRequest)
		return
	}

	foundTask, ok := s.tasks[id]
	if !ok {
		s.log.Warn("task not found")
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(&foundTask); err != nil {
		s.log.Error("failed to marshal single task JSON", "id", id, "error", err, "remote_addr", r.RemoteAddr)
		http.Error(w, "failed to marshal single task JSON", http.StatusInternalServerError)
		return
	}
}

func (s *Server) CreateTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task

	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		s.log.Warn("failed to decode task JSON", "error", err, "remote_addr", r.RemoteAddr)
		http.Error(w, "bad request body", http.StatusBadRequest)
		return
	}

	if len(newTask.Title) > 200 || len(newTask.Description) > 2000 { // bad practice. Business logic should not be in HTTP-handler.
		s.log.Warn("client try send to long data")
		http.Error(w, "bad request body", http.StatusBadRequest)
		return
	}

	if newTask.Title == "" { // bad practice. Business logic should not be in HTTP-handler.
		s.log.Warn("client sent empty title", "title =", newTask.Title)
		newTask.Title = "Default title"
	}

	newTask.ID = s.counter.Add(1)
	s.tasks[newTask.ID] = newTask

	w.Header().Set("Content-Type", "application/json")

	payload, err := json.Marshal(newTask)
	if err != nil {
		s.log.Error("failed to encode new task", "error", err, "remote_addr", r.RemoteAddr)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(payload)
}

func (s *Server) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		s.log.Warn("invalid id format")
		http.Error(w, "invalid id format", http.StatusBadRequest)
		return
	}

	foundTask, ok := s.tasks[id]
	if !ok {
		s.log.Warn("task not found")
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&foundTask); err != nil {
		s.log.Warn("failed to decode task JSON", "error", err, "remote_addr", r.RemoteAddr)
		http.Error(w, "bad request body", http.StatusBadRequest)
		return
	}

	foundTask.ID = id
	s.tasks[id] = foundTask

	if err := json.NewEncoder(w).Encode(&foundTask); err != nil {
		s.log.Error("failed to encode updated task", "error", err, "remote_addr", r.RemoteAddr)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (s *Server) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		s.log.Warn("invalid id format")
		http.Error(w, "invalid id format", http.StatusBadRequest)
		return
	}

	_, ok := s.tasks[id]
	if !ok {
		s.log.Warn("task not found")
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	delete(s.tasks, id)

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	tasks := map[int64]Task{
		0: {ID: 0, Title: "Write simple HTTP API", Description: "This is quite difficult for me", Done: false},
		1: {ID: 1, Title: "Drink tea and keep calm", Description: "My english is very the capital of GB (jk)", Done: false},
		2: {ID: 2, Title: "Fix any bugs", Done: false},
	}

	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))

	s := NewServer(tasks, l)
	mux := http.NewServeMux()

	mux.HandleFunc("GET /tasks", s.ListTasks)
	mux.HandleFunc("POST /tasks", s.CreateTask)
	mux.HandleFunc("GET /tasks/{id}", s.GetTaskByID)
	mux.HandleFunc("PUT /tasks/{id}", s.UpdateTask)
	mux.HandleFunc("DELETE /tasks/{id}", s.DeleteTask)

	s.log.Info("http server started", "port", 8080)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		s.log.Error("HTTP server stopped", "error", err)
	}
}
