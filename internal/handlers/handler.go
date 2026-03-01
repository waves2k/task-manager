package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/waves2k/task-manager/internal/service"
)

type Handler struct {
	todoService service.TodoService
}

func NewHandler(todoService service.TodoService) *Handler {
	return &Handler{
		todoService: todoService,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.SetTrustedProxies(nil)

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api")
	{
		todo := api.Group("/todo")
		{
			todo.POST("/")
		}
		api.GET("/default", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"message": "TODO Api is running!",
				"status":  "success",
			})
		})

	}

	return router
}
