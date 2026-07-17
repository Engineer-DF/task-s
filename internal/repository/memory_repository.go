package repository

import (
	"errors"
	"fmt"
	"slices"
	"sync/atomic"

	"github.com/Engineer-DF/task-s/internal/domain"
)

var ErrNotFound error = errors.New("repository: not found")

type TaskRepository struct {
	tasks   map[int64]domain.Task
	counter atomic.Int64
}

func NewTaskRepository(initialTasks map[int64]domain.Task) *TaskRepository {
	repo := &TaskRepository{
		tasks: initialTasks,
	}
	repo.counter.Store(repo.findMaxID())
	return repo
}

func (r *TaskRepository) findMaxID() int64 {
	var maxID int64 = -1

	for key := range r.tasks {
		if key > maxID {
			maxID = key
		}
	}
	return maxID
}

func (r *TaskRepository) GetAll() ([]domain.Task, error) {
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

func (r *TaskRepository) GetByID(id int64) (domain.Task, error) {
	task, exists := r.tasks[id]
	if !exists {
		return domain.Task{}, fmt.Errorf("failed to get task (ID: %d): %w", id, ErrNotFound)
	}

	return task, nil
}

// Возвращаем ошибку для возможной реализации метода с другими видами хранения данных (БД),
// где требуется корректно обработать ошибку.
func (r *TaskRepository) Create(task domain.Task) (domain.Task, error) {
	task.ID = r.counter.Add(1)
	r.tasks[task.ID] = task

	return task, nil
}

func (r *TaskRepository) Update(id int64, task domain.Task) error {
	_, exists := r.tasks[id]
	if !exists {
		return fmt.Errorf("failed to update task (ID: %d): %w", id, ErrNotFound)
	}

	task.ID = id
	r.tasks[id] = task

	return nil
}

func (r *TaskRepository) Delete(id int64) error {
	_, exists := r.tasks[id]
	if !exists {
		return fmt.Errorf("failed to delete task (ID: %d): %w", id, ErrNotFound)
	}
	delete(r.tasks, id)
	return nil
}
