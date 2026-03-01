package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/waves2k/task-manager/internal/database"
	"github.com/waves2k/task-manager/internal/models"
)

type TodoRepository interface {
	CreateTodo(ctx context.Context, title string, completed bool) (*models.Todo, error)
}

type todoRepository struct {
	pool *pgxpool.Pool
}

func NewTodoService(pool *pgxpool.Pool) *todoRepository {
	return &todoRepository{
		pool: pool,
	}
}

func (r *todoRepository) CreateTodo(ctx context.Context, title string, completed bool) (*models.Todo, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	query := fmt.Sprintf(`
		INSERT INTO %s (title, completed)
		VALUES ($1, $2)
		RETURNING id, title, completed, created_at, updated_at
	`, database.TodoTable)

	var todo models.Todo
	if err := r.pool.QueryRow(ctx, query, title, completed).Scan(
		&todo.Id,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &todo, nil
}
