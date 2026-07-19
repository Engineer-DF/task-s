package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Engineer-DF/task-s/internal/domain"
	"github.com/Engineer-DF/task-s/internal/service"
)

// Этот слой НЕ должен видеть никакие другие слои кроме domain и service
type TaskHandler struct {
	service *service.Service
	log     *slog.Logger
}

func NewTaskHandler(service *service.Service, log *slog.Logger) *TaskHandler {
	return &TaskHandler{
		service: service,
		log:     log,
	}
}

func (t *TaskHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	tasks, err := t.service.GetAll()
	if err != nil {
		t.log.Error("failed to get all tasks", "method", r.Method, "pattern", r.Pattern, "error", err)
		http.Error(w, "failed to get all tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		t.log.Error("failed to encode JSON", "method", r.Method, "pattern", r.Pattern, "error", err)
		http.Error(w, "failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}

func (t *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		t.log.Warn("invalid ID format")
		http.Error(w, "invalid ID format", http.StatusBadRequest)
		return
	}

	foundTask, err := t.service.GetByID(id)
	if err != nil {
		t.log.Warn("task not found")
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(foundTask); err != nil {
		t.log.Error("failed to encode JSON", "method", r.Method, "pattern", r.Pattern, "error", err)
		http.Error(w, "failed to marshal JSON", http.StatusInternalServerError)
		return
	}
}

func (t *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var gotTask domain.Task

	if err := json.NewDecoder(r.Body).Decode(&gotTask); err != nil {
		t.log.Warn("failed to decode JSON", "method", r.Method, "pattern", r.Pattern, "error", err)
		http.Error(w, "bad request body", http.StatusBadRequest)
		return
	}

	newTask, err := t.service.Create(gotTask)
	if err != nil {
		t.log.Error("failed to create task", "method", r.Method, "pattern", r.Pattern, "error", err)
		http.Error(w, "failed to create task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(&newTask); err != nil {
		t.log.Error("failed to encode JSON", "method", r.Method, "pattern", r.Pattern, "error", err)
		http.Error(w, "failed to marshal JSON", http.StatusInternalServerError)
		return
	}
}

func (t *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	var gotTask domain.Task

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		t.log.Warn("invalid ID format")
		http.Error(w, "invalid ID format", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&gotTask); err != nil {
		t.log.Warn("failed to decode JSON", "method", r.Method, "pattern", r.Pattern, "error", err)
		http.Error(w, "bad request body", http.StatusBadRequest)
		return
	}

	err = t.service.Update(id, gotTask)
	if err != nil {
		t.log.Warn("task not found")
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(&gotTask); err != nil {
		t.log.Error("failed to encode JSON", "method", r.Method, "pattern", r.Pattern, "error", err)
		http.Error(w, "failed to marshal JSON", http.StatusInternalServerError)
		return
	}
}

func (t *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		t.log.Warn("invalid ID format")
		http.Error(w, "invalid ID format", http.StatusBadRequest)
		return
	}

	err = t.service.Delete(id)
	if err != nil {
		t.log.Warn("task not found")
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
