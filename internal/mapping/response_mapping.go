package mapping

import (
	"github.com/waves2k/task-manager/internal/contracts"
	"github.com/waves2k/task-manager/internal/models"
)

func MapToTodoResponse(model *models.Todo) *contracts.TodoResponse {
	return &contracts.TodoResponse{
		Id:        model.Id,
		Title:     model.Title,
		Completed: model.Completed,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

func MapToUserResponse(model *models.User) *contracts.UserResponse {
	return &contracts.UserResponse{
		Id:        model.Id,
		Email:     model.Email,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

func MapToListResponse(model *models.List) contracts.ListResponse {
	return contracts.ListResponse{
		Id:          model.Id,
		Title:       model.Title,
		Description: model.Description,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}

func MapToListsSliceReponse(lists []models.List) []contracts.ListResponse {
	mappedLists := make([]contracts.ListResponse, len(lists))
	for i := 0; i < len(lists); i++ {
		mappedLists[i] = MapToListResponse(&lists[i])
	}
	return mappedLists
}
