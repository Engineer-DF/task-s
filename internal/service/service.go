package service

import (
	"context"
	"errors"

	"github.com/Engineer-DF/task-s/internal/domain"
)

var (
	ErrInvalidTask = errors.New("invalid task")
	ErrTooLongTask = errors.New("task title/description too long")
)

// Интерфейс существует для того, чтобы отделить реализацию методов от получателя
// Интерфейс обычно объявляет потребитель, а не тот, кто его реализует.
type Repository interface {
	GetAll(ctx context.Context) ([]domain.Task, error)
	GetByID(ctx context.Context, id int64) (domain.Task, error)
	Create(ctx context.Context, task domain.Task) (domain.Task, error)
	Update(ctx context.Context, id int64, task domain.Task) error
	Delete(ctx context.Context, id int64) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
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

func (s *Service) GetAll(ctx context.Context) ([]domain.Task, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) GetByID(ctx context.Context, id int64) (domain.Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, task domain.Task) (domain.Task, error) {
	if err := ctx.Err(); err != nil {
		return domain.Task{}, err
	}

	err := s.validateTask(task)

	if err != nil {
		return domain.Task{}, err
	}

	return s.repo.Create(ctx, task)
}

func (s *Service) Update(ctx context.Context, id int64, task domain.Task) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	err := s.validateTask(task)

	if err != nil {
		return err
	}

	return s.repo.Update(ctx, id, task)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
