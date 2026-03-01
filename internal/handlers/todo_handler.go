package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/waves2k/task-manager/internal/contracts"
	"github.com/waves2k/task-manager/internal/service"
	"github.com/waves2k/task-manager/internal/utils/mapper"
)

func (h *Handler) updateTodo() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		unparsedId, err := strconv.Atoi(idParam)
		if err != nil {
			NewErrorResponse(c, http.StatusBadRequest, InvalidParamData)
			return
		}
		var input contracts.UpdateTodoRequest
		if err := c.BindJSON(&input); err != nil {
			NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		if err := input.Validate(); err != nil {
			NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		todo, err := h.todoService.Update(c, service.UpdateTodoParams{
			Id:        unparsedId,
			Title:     input.Title,
			Completed: input.Completed,
		})
		if err != nil {
			if errors.Is(sql.ErrNoRows, err) {
				NewErrorResponse(c, http.StatusNotFound, err.Error())
				return
			}
			NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		NewSucceededResponse(c, http.StatusOK, mapper.MapToTodoResponse(todo))
	}
}

func (h *Handler) getTodoById() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		unparsedId, err := strconv.Atoi(idParam)
		if err != nil {
			NewErrorResponse(c, http.StatusBadRequest, InvalidParamData)
			return
		}
		todo, err := h.todoService.GetById(c, unparsedId)
		if err != nil {
			if errors.Is(sql.ErrNoRows, err) {
				NewErrorResponse(c, http.StatusNotFound, err.Error())
				return
			}
			NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		NewSucceededResponse(c, http.StatusOK, mapper.MapToTodoResponse(todo))
	}
}

func (h *Handler) createTodoHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input contracts.CreateTodoRequest
		if err := c.ShouldBindJSON(&input); err != nil {
			NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		// TODO: GET USER_ID FROM THE CONTEXT
		todo, err := h.todoService.Create(c.Request.Context(), input.Title, input.Completed)
		if err != nil {
			NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		NewSucceededResponse(c, http.StatusCreated, mapper.MapToTodoResponse(todo))
	}
}

func (h *Handler) deleteTodo() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		unparsedId, err := strconv.Atoi(idParam)
		if err != nil {
			NewErrorResponse(c, http.StatusBadRequest, InvalidParamData)
			return
		}
		deletedId, err := h.todoService.Delete(c, unparsedId)
		if err != nil {
			if errors.Is(sql.ErrNoRows, err) {
				NewErrorResponse(c, http.StatusNotFound, err.Error())
				return
			}
			NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		NewSucceededResponse(c, http.StatusCreated, map[string]interface{}{
			"id": deletedId,
		})
	}
}
