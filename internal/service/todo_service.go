package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	apperrors "github.com/waves2k/task-manager/internal/errors"
	"github.com/waves2k/task-manager/internal/models"
	"github.com/waves2k/task-manager/internal/repository"
)

type TodoService interface {
	Create(ctx context.Context, title string, completed bool, listId, userId uuid.UUID) (*models.Todo, error)
	GetAll(ctx context.Context, listId, userId uuid.UUID) ([]models.Todo, error)
	GetById(ctx context.Context, todoId uuid.UUID) (*models.Todo, error)
	Update(ctx context.Context, params UpdateTodoParams) (*models.Todo, error)
	Delete(ctx context.Context, todoId uuid.UUID) (uuid.UUID, error)
}

type todoService struct {
	todoRepo repository.TodoRepository
	listRepo repository.ListRepository
	userRepo repository.UserRepository
}

func CreateNewTodoService(
	todoRepo repository.TodoRepository,
	listRepo repository.ListRepository,
	userRepo repository.UserRepository,
) *todoService {
	return &todoService{
		todoRepo: todoRepo,
		listRepo: listRepo,
		userRepo: userRepo,
	}
}

type UpdateTodoParams struct {
	Id        uuid.UUID
	Title     *string
	Completed *bool
	UserId    uuid.UUID
}

func (s *todoService) Create(ctx context.Context, title string, completed bool, listId, userId uuid.UUID) (*models.Todo, error) {
	if !s.userRepo.IsExists(ctx, userId) {
		return nil, apperrors.NewNotFoundError(apperrors.UserNotFoundMessage, fmt.Errorf(apperrors.UserNotFoundMessage))
	}
	if !s.listRepo.IsExists(ctx, listId) {
		return nil, apperrors.NewNotFoundError(apperrors.ListNotFoundMessage, fmt.Errorf(apperrors.ListNotFoundMessage))
	}
	return s.todoRepo.Create(ctx, title, completed, listId)
}

func (s *todoService) GetAll(ctx context.Context, listId, userId uuid.UUID) ([]models.Todo, error) {
	return s.todoRepo.GetAll(ctx, listId, userId)
}

func (s *todoService) GetById(ctx context.Context, todoId uuid.UUID) (*models.Todo, error) {
	return s.todoRepo.GetById(ctx, todoId)
}

func (s *todoService) Update(ctx context.Context, params UpdateTodoParams) (*models.Todo, error) {
	return s.todoRepo.Update(ctx, repository.UpdateTodoParams{
		Id:        params.Id,
		Title:     params.Title,
		Completed: params.Completed,
	})
}

func (s *todoService) Delete(ctx context.Context, todoId uuid.UUID) (uuid.UUID, error) {
	if err := s.todoRepo.Delete(ctx, todoId); err != nil {
		return uuid.Nil, err
	}
	return todoId, nil
}
