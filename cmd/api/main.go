package main

import (
	"log"
	"net/http"
	"time"

	_ "github.com/joho/godotenv"
	"github.com/waves2k/task-manager/internal/config"
	"github.com/waves2k/task-manager/internal/database"
	"github.com/waves2k/task-manager/internal/handlers"
	"github.com/waves2k/task-manager/internal/repository"
	"github.com/waves2k/task-manager/internal/service"
)

func main() {

	// Loading the configuration.
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Configuring the database.
	pool, err := database.Connect(cfg.ConnectionString)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer pool.Close()

	// Configuring dependences.
	todoRepo := repository.NewTodoService(pool)
	todoService := service.CreateNewTodoService(todoRepo)
	handler := handlers.NewHandler(todoService)

	// Configuring the server.
	// Вынести в отдельную сущность в идеале
	httpServer := http.Server{
		Addr:           ":" + cfg.Port,
		Handler:        handler.InitRoutes(),
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	// Listening to requests.
	httpServer.ListenAndServe()
}
