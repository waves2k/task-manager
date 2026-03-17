package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/waves2k/task-manager/internal/api"
	"github.com/waves2k/task-manager/internal/contracts"
	"github.com/waves2k/task-manager/internal/mapping"
	"github.com/waves2k/task-manager/internal/service"
)

func (h *Handler) createList() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input contracts.CreateListRequest
		if err := TryBindRequestBody(c, &input); err != nil {
			c.Error(err)
			return
		}
		userId, err := GetUserId(c)
		if err != nil {
			c.Error(err)
			return
		}
		list, err := h.listSErvice.Create(c.Request.Context(), input.Title, input.Description, userId)
		if err != nil {
			c.Error(err)
			return
		}
		api.NewSucceededResponse(c, http.StatusCreated, mapping.MapToListResponse(list))
	}
}

func (h *Handler) getAllLists() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := GetUserId(c)
		if err != nil {
			c.Error(err)
			return
		}
		lists, err := h.listSErvice.GetAll(c.Request.Context(), userId)
		if err != nil {
			c.Error(err)
			return
		}
		api.NewSucceededResponse(c, http.StatusCreated, mapping.MapToListsSliceReponse(lists))
	}
}

func (h *Handler) getListById() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := GetUserId(c)
		if err != nil {
			c.Error(err)
			return
		}
		listId, err := GetListId(c)
		if err != nil {
			c.Error(err)
			return
		}
		list, err := h.listSErvice.GetById(c.Request.Context(), listId, userId)
		if err != nil {
			c.Error(err)
			return
		}
		api.NewSucceededResponse(c, http.StatusCreated, mapping.MapToListResponse(list))
	}
}

func (h *Handler) updateList() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := GetUserId(c)
		if err != nil {
			c.Error(err)
			return
		}
		listId, err := GetListId(c)
		if err != nil {
			c.Error(err)
			return
		}
		var input contracts.UpdateListRequest
		if err := TryBindRequestBody(c, &input); err != nil {
			c.Error(err)
			return
		}
		if err = input.Validate(); err != nil {
			c.Error(err)
			return
		}
		list, err := h.listSErvice.Update(c.Request.Context(), service.UpdateListParams{
			Id:          listId,
			Title:       input.Title,
			Description: input.Description,
			UserId:      userId,
		})
		if err != nil {
			c.Error(err)
			return
		}
		api.NewSucceededResponse(c, http.StatusCreated, mapping.MapToListResponse(list))
	}
}

func (h *Handler) deleteList() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := GetUserId(c)
		if err != nil {
			c.Error(err)
			return
		}
		listId, err := GetListId(c)
		if err != nil {
			c.Error(err)
			return
		}
		list, err := h.listSErvice.GetById(c.Request.Context(), listId, userId)
		if err != nil {
			c.Error(err)
			return
		}
		api.NewSucceededResponse(c, http.StatusCreated, mapping.MapToListResponse(list))
	}

}
