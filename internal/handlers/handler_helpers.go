package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	apperrors "github.com/waves2k/task-manager/internal/errors"
)

func GetUserId(c *gin.Context) (uuid.UUID, error) {
	unparsed, ok := c.Get("user_id")
	if !ok {
		return uuid.Nil, apperrors.NewValidationError(apperrors.UnauthorizedErrorMessage, fmt.Errorf(apperrors.ValidationErrorMessage))
	}
	userId, ok := unparsed.(uuid.UUID)
	if !ok {
		return uuid.Nil, apperrors.NewValidationError(apperrors.UnauthorizedErrorMessage, fmt.Errorf(apperrors.ValidationErrorMessage))
	}
	return userId, nil
}

func GetListId(c *gin.Context) (uuid.UUID, error) {
	listId, err := uuid.Parse(c.Param("listId"))
	if err != nil {
		return uuid.Nil, apperrors.NewWrongInputDataError(apperrors.WrongParamsInputData, err)
	}
	return listId, nil
}

func TryBindRequestBody(c *gin.Context, input interface{}) error {
	if err := c.BindJSON(&input); err != nil {
		return apperrors.NewWrongInputDataError(apperrors.InvalidRequestBodyData, err)
	}
	return nil
}
