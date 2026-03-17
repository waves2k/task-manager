package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/waves2k/task-manager/internal/api"
	"github.com/waves2k/task-manager/internal/config"
	"github.com/waves2k/task-manager/internal/database"
	"github.com/waves2k/task-manager/internal/handlers"
	"github.com/waves2k/task-manager/internal/middleware"
	"github.com/waves2k/task-manager/internal/repository"
	"github.com/waves2k/task-manager/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	pool, err := database.Connect(cfg.ConnectionString)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer pool.Close()

	todoRepo := repository.NewTodoService(pool)
	userRepo := repository.NewUserRepository(pool)
	listRepo := repository.NewListRepository(pool)

	passwordHasher := service.NewPasswordHasher(cfg.PasswordSalt)
	tokenService := service.NewJWTService(cfg.JWTSecret)

	userService := service.NewUserService(userRepo, passwordHasher, tokenService)
	todoService := service.CreateNewTodoService(todoRepo, listRepo, userRepo)
	listService := service.NewListService(listRepo)

	authMiddleware := middleware.NewAuthMiddleware(tokenService)
	errorHandelerMiddleware := middleware.NewErrorHandlerMiddleware()

	handler := handlers.NewHandler(
		todoService,
		userService,
		listService,
		authMiddleware,
		errorHandelerMiddleware,
	)

	srv := api.CreateServer(cfg.Port, handler.InitRoutes())

	go func() {
		if err := srv.Run(); err != nil {
			log.Fatal("Error during server listening:", err)
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-done

	logrus.Print("Shutting down the server..")
}
