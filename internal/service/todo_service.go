package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/waves2k/task-manager/internal/models"
	"github.com/waves2k/task-manager/internal/repository"
)

type TodoService interface {
	Create(ctx context.Context, title string, completed bool) (*models.Todo, error)
	GetAll(ctx context.Context) (*[]models.Todo, error)
	GetById(ctx context.Context, todoId int) (*models.Todo, error)
	Update(ctx context.Context, params UpdateTodoParams) (*models.Todo, error)
	Delete(ctx context.Context, todoId int) (int, error)
}

type todoService struct {
	todoRepo repository.TodoRepository
}

func CreateNewTodoService(todoRepo repository.TodoRepository) *todoService {
	return &todoService{
		todoRepo: todoRepo,
	}
}

type UpdateTodoParams struct {
	Id        int
	Title     *string
	Completed *bool
	UserId    uuid.UUID
}

func (s *todoService) Create(ctx context.Context, title string, completed bool) (*models.Todo, error) {
	return s.todoRepo.Create(ctx, title, completed)
}

func (s *todoService) GetAll(ctx context.Context) (*[]models.Todo, error) {
	return nil, nil
}

func (s *todoService) GetById(ctx context.Context, todoId int) (*models.Todo, error) {
	return s.todoRepo.GetById(ctx, todoId)
}

func (s *todoService) Update(ctx context.Context, params UpdateTodoParams) (*models.Todo, error) {
	return s.todoRepo.Update(ctx, repository.UpdateTodoParams{
		Id:        params.Id,
		Title:     params.Title,
		Completed: params.Completed,
	})
}

func (s *todoService) Delete(ctx context.Context, todoId int) (int, error) {
	if err := s.todoRepo.Delete(ctx, todoId); err != nil {
		return 0, err
	}
	return todoId, nil
}
