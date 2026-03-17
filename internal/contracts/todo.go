package contracts

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	apperrors "github.com/waves2k/task-manager/internal/errors"
)

type CreateTodoRequest struct {
	Title     string `json:"title" binding:"required"`
	Completed bool   `json:"completed"`
}

type TodoResponse struct {
	Id        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateTodoRequest struct {
	Title     *string `json:"title"`
	Completed *bool   `json:"completed"`
}

func (i *UpdateTodoRequest) Validate() error {
	if i.Completed == nil &&
		i.Title == nil {
		return apperrors.NewValidationError(apperrors.InvalidRequestBodyData, fmt.Errorf(apperrors.InvalidRequestBodyData))
	}
	return nil
}
