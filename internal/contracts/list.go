package contracts

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	apperrors "github.com/waves2k/task-manager/internal/errors"
)

type CreateListRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type ListResponse struct {
	Id          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UpdateListRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description" `
}

func (i *UpdateListRequest) Validate() error {
	if i.Title == nil &&
		i.Description == nil {
		return apperrors.NewValidationError(apperrors.InvalidRequestBodyData, fmt.Errorf(apperrors.InvalidRequestBodyData))
	}
	return nil
}
