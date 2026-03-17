package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/waves2k/task-manager/internal/models"
	"github.com/waves2k/task-manager/internal/repository"
)

// INJECT TODO REPO INTO LIST_SERVICE
// MAKE TWO TYPES OF GET METHOD(WITH TODOS AND WITHOUT)
type ListService interface {
	Create(ctx context.Context, title, description string, userId uuid.UUID) (*models.List, error)
	GetAll(ctx context.Context, userId uuid.UUID) ([]models.List, error)
	GetById(ctx context.Context, listId, userId uuid.UUID) (*models.List, error)
	Update(ctx context.Context, params UpdateListParams) (*models.List, error)
	Delete(ctx context.Context, listId, userId uuid.UUID) (uuid.UUID, error)
}

type listService struct {
	listRepo repository.ListRepository
}

func NewListService(listRepo repository.ListRepository) *listService {
	return &listService{
		listRepo: listRepo,
	}
}

type UpdateListParams struct {
	Id          uuid.UUID
	Title       *string
	Description *string
	UserId      uuid.UUID
}

func (s *listService) Create(ctx context.Context, title, description string, userId uuid.UUID) (*models.List, error) {
	list, err := s.listRepo.Create(ctx, title, description, userId)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s *listService) GetAll(ctx context.Context, userId uuid.UUID) ([]models.List, error) {
	lists, err := s.listRepo.GetAll(ctx, userId)
	if err != nil {
		return nil, err
	}
	return lists, err
}

func (s *listService) GetById(ctx context.Context, listId, userId uuid.UUID) (*models.List, error) {
	list, err := s.listRepo.GetById(ctx, listId, userId)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s *listService) Update(ctx context.Context, params UpdateListParams) (*models.List, error) {
	list, err := s.listRepo.Update(ctx, repository.UpdateListParams{
		Id:          params.Id,
		Title:       params.Title,
		Description: params.Description,
		UserId:      params.UserId,
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s *listService) Delete(ctx context.Context, listId, userId uuid.UUID) (uuid.UUID, error) {
	if err := s.listRepo.Delete(ctx, listId, userId); err != nil {
		return uuid.Nil, err
	}
	return listId, nil
}
