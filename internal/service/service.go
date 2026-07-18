package service

import (
	"errors"

	"github.com/Engineer-DF/task-s/internal/domain"
)

var (
	ErrInvalidTask = errors.New("invalid task")
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

func (s *Service) validateTask(task domain.Task) error {
	if len(task.Title) > 100 || len(task.Description) > 300 {
		return ErrTooLongTask
	}

	if task.Title == "" {
		return ErrInvalidTask
	}

	return nil
}

func (s *Service) GetAll() ([]domain.Task, error) {
	return s.repo.GetAll()
}

func (s *Service) GetByID(id int64) (domain.Task, error) {
	return s.repo.GetByID(id)
}

func (s *Service) Create(task domain.Task) (domain.Task, error) {
	err := s.validateTask(task)

	if err != nil {
		return domain.Task{}, err
	}

	return s.repo.Create(task)
}

func (s *Service) Update(id int64, task domain.Task) error {
	err := s.validateTask(task)

	if err != nil {
		return err
	}

	return s.repo.Update(id, task)
}

func (s *Service) Delete(id int64) error {
	return s.repo.Delete(id)
}
