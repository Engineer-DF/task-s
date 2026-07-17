package service

import (
	"errors"

	"github.com/Engineer-DF/task-s/internal/domain"
	"github.com/Engineer-DF/task-s/internal/repository"
)

var (
	ErrInvalidTask = errors.New("task not found")
	ErrTooLongTask = errors.New("task title/description too long")
)

// Интерфейс существует для того, чтобы отделить реализацию методов от получателя
type Repository interface {
	GetAll() ([]domain.Task, error)
	GetByID(id int64) (domain.Task, error)
	Create(task domain.Task) (domain.Task, error)
	Update(id int64, task domain.Task) error
	Delete(id int64) error
}

type Service struct {
	repo Repository
}

// TODO: написать crud методы для реализации интерфейса

func validateTask(task domain.Task, err error) error {
	if len(task.Title) > 100 || len(task.Description) > 300 {
		return ErrTooLongTask
	}

	if errors.Is(err, repository.ErrNotFound) { // перенести в crud методы
		return ErrInvalidTask
	}

	return nil
}
