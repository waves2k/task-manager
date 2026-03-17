package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/waves2k/task-manager/internal/middleware"
	"github.com/waves2k/task-manager/internal/service"
)

type Handler struct {
	authMiddleware *middleware.AuthMiddleware
	errorHandler   *middleware.ErrorHandlerMiddleware
	todoService    service.TodoService
	userService    service.UserService
	listSErvice    service.ListService
}

func NewHandler(
	todoService service.TodoService,
	userService service.UserService,
	listSErvice service.ListService,
	authMiddleware *middleware.AuthMiddleware,
	errorHandler *middleware.ErrorHandlerMiddleware,
) *Handler {
	return &Handler{
		todoService:    todoService,
		userService:    userService,
		listSErvice:    listSErvice,
		authMiddleware: authMiddleware,
		errorHandler:   errorHandler,
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

	api := router.Group("/api", h.authMiddleware.Handler(), h.errorHandler.Handler())
	{
		list := api.Group("/list")
		{
			list.POST("/", h.createList())
			list.GET("/", h.getAllLists())
			list.GET("/:id", h.getListById())
			list.PUT("/:id", h.updateList())
			list.DELETE("/:id", h.deleteList())

			todo := list.Group("/lists/:id/todo")
			{
				todo.POST("/", h.createTodoHandler())
				todo.GET("/", h.createTodoHandler())
			}
		}

		todos := api.Group("/todo")
		{
			todos.PUT("/:id", h.updateTodo())
			todos.GET("/:id", h.getTodoById())
			todos.DELETE("/:id", h.deleteTodo())
		}
		api.GET("/default", func(c *gin.Context) {
			authHeader := c.GetHeader("Authorization")

			c.JSON(200, gin.H{
				"message": "TODO Api is running!",
				"status":  "success",
				"header":  authHeader,
			})
		})

	}

	return router
}
