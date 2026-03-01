package mapper

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
