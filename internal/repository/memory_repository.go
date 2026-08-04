package repository

import (
	"context"
	"fmt"
	"slices"
	"sync"

	"github.com/Engineer-DF/task-s/internal/domain"
)

type TaskRepository struct {
	mu      sync.RWMutex
	tasks   map[int64]domain.Task
	counter int64
}

func NewTaskRepository(initialTasks map[int64]domain.Task) *TaskRepository {
	copiedTasks := make(map[int64]domain.Task, len(initialTasks))
	var maxID int64 = -1

	for key, value := range initialTasks {
		copiedTasks[key] = value
		if key > maxID {
			maxID = key
		}
	}

	repo := &TaskRepository{
		tasks:   copiedTasks,
		counter: maxID,
	}
	return repo
}

/*func (r *TaskRepository) findMaxID() int64 {
	var maxID int64 = -1

	r.mu.Lock()
	defer r.mu.Unlock()

	for key := range r.tasks {
		if key > maxID {
			maxID = key
		}
	}
	return maxID
}*/

func (r *TaskRepository) GetAll(ctx context.Context) ([]domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	mapKeys := make([]int64, 0, len(r.tasks))

	for keys := range r.tasks {
		mapKeys = append(mapKeys, keys)
	}
	slices.Sort(mapKeys)

	allTasks := make([]domain.Task, 0, len(r.tasks))

	for _, keys := range mapKeys {
		allTasks = append(allTasks, r.tasks[keys])
	}

	return allTasks, nil
}

func (r *TaskRepository) GetByID(ctx context.Context, id int64) (domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, exists := r.tasks[id]
	if !exists {
		return domain.Task{}, fmt.Errorf("in-memory map get task failed (ID: %d): %w", id, domain.ErrNotFound)
	}

	return task, nil
}

// Возвращаем ошибку для возможной реализации метода с другими видами хранения данных (БД),
// где требуется корректно обработать ошибку.
func (r *TaskRepository) Create(ctx context.Context, task domain.Task) (domain.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.counter++
	task.ID = r.counter
	r.tasks[task.ID] = task

	return task, nil
}

func (r *TaskRepository) Update(ctx context.Context, id int64, task domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.tasks[id]
	if !exists {
		return fmt.Errorf("in-memory map update task failed (ID: %d): %w", id, domain.ErrNotFound)
	}

	task.ID = id
	r.tasks[id] = task

	return nil
}

func (r *TaskRepository) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.tasks[id]
	if !exists {
		return fmt.Errorf("in-memory map delete task failed (ID: %d): %w", id, domain.ErrNotFound)
	}
	delete(r.tasks, id)

	return nil
}
