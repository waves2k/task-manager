package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/waves2k/task-manager/internal/api"
	"github.com/waves2k/task-manager/internal/contracts"
	apperrors "github.com/waves2k/task-manager/internal/errors"
	"github.com/waves2k/task-manager/internal/mapping"
)

func (h *Handler) signUp(c *gin.Context) {
	var input contracts.CreateUserRequest
	if err := TryBindRequestBody(c, &input); err != nil {
		c.Error(err)
		return
	}
	if len(input.Password) < 6 {
		c.Error(apperrors.NewValidationError(apperrors.ShortPasswordErrorData, fmt.Errorf(apperrors.ShortPasswordErrorData)))
		return
	}
	user, err := h.userService.Register(c.Request.Context(), input.Email, input.Password)
	if err != nil {
		c.Error(err)
		return
	}
	api.NewSucceededResponse(c, http.StatusCreated, mapping.MapToUserResponse(user))
}

func (h *Handler) signIn(c *gin.Context) {
	var input contracts.LoginUserRequest
	if err := TryBindRequestBody(c, &input); err != nil {
		c.Error(err)
		return
	}
	token, err := h.userService.Login(c.Request.Context(), input.Email, input.Password)
	if err != nil {
		c.Error(err)
		return
	}
	api.NewCustomDataSucceededResponse(c, http.StatusCreated, map[string]interface{}{
		"token": token})
}
