package handlers

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) signUp(c *gin.Context) {
	c.JSON(200, "Reg")
}

func (h *Handler) signIn(c *gin.Context) {
	c.JSON(200, "Auth")
}
