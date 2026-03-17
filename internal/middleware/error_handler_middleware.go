package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/waves2k/task-manager/internal/api"
	apperrors "github.com/waves2k/task-manager/internal/errors"
)

type ErrorHandlerMiddleware struct {
}

func NewErrorHandlerMiddleware() *ErrorHandlerMiddleware {
	return &ErrorHandlerMiddleware{}
}

func (m *ErrorHandlerMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		logrus.Errorf("Error: %v", err)

		switch {
		case apperrors.IsValidation(err) || apperrors.IsWrongInputData(err):
			api.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		case apperrors.IsFrobidden(err):
			api.NewErrorResponse(c, http.StatusForbidden, err.Error())
		case apperrors.IsUnauthorized(err):
			api.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		case apperrors.IsNotFound(err):
			api.NewErrorResponse(c, http.StatusNotFound, err.Error())
		case apperrors.IsAlreadyExists(err):
			api.NewErrorResponse(c, http.StatusConflict, err.Error())
		case apperrors.IsInternal(err):
			api.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		}

		c.Abort()
	}
}
