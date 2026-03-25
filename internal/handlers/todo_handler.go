package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/waves2k/task-manager/internal/api"
	"github.com/waves2k/task-manager/internal/contracts"
	"github.com/waves2k/task-manager/internal/mapping"
	"github.com/waves2k/task-manager/internal/service"
)

func (h *Handler) updateTodo() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")

		unparsedId, err := uuid.Parse(idParam)
		if err != nil {
			api.NewErrorResponse(c, http.StatusBadRequest, api.InvalidParamData)
			return
		}
		var input contracts.UpdateTodoRequest
		if err := c.BindJSON(&input); err != nil {
			api.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		if err := input.Validate(); err != nil {
			api.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		todo, err := h.todoService.Update(c, service.UpdateTodoParams{
			Id:        unparsedId,
			Title:     input.Title,
			Completed: input.Completed,
		})
		if err != nil {
			if errors.Is(sql.ErrNoRows, err) {
				api.NewErrorResponse(c, http.StatusNotFound, err.Error())
				return
			}
			api.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		api.NewSucceededResponse(c, http.StatusOK, mapping.MapToTodoResponse(todo))
	}
}

func (h *Handler) getTodoById() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		unparsedId, err := uuid.Parse(idParam)
		if err != nil {
			api.NewErrorResponse(c, http.StatusBadRequest, api.InvalidParamData)
			return
		}
		todo, err := h.todoService.GetById(c, unparsedId)
		if err != nil {
			if errors.Is(sql.ErrNoRows, err) {
				api.NewErrorResponse(c, http.StatusNotFound, err.Error())
				return
			}
			api.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		api.NewSucceededResponse(c, http.StatusOK, mapping.MapToTodoResponse(todo))
	}
}

func (h *Handler) createTodoHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input contracts.CreateTodoRequest
		if err := c.ShouldBindJSON(&input); err != nil {
			api.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		// TODO: GET USER_ID FROM THE CONTEXT
		todo, err := h.todoService.Create(c.Request.Context(), input.Title, input.Completed, input.ListId)
		if err != nil {
			api.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		api.NewSucceededResponse(c, http.StatusCreated, mapping.MapToTodoResponse(todo))
	}
}

func (h *Handler) deleteTodo() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		unparsedId, err := uuid.Parse(idParam)
		if err != nil {
			api.NewErrorResponse(c, http.StatusBadRequest, api.InvalidParamData)
			return
		}
		deletedId, err := h.todoService.Delete(c, unparsedId)
		if err != nil {
			if errors.Is(sql.ErrNoRows, err) {
				api.NewErrorResponse(c, http.StatusNotFound, err.Error())
				return
			}
			api.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		api.NewSucceededResponse(c, http.StatusCreated, map[string]interface{}{
			"id": deletedId,
		})
	}
}
