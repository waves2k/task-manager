package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/waves2k/task-manager/internal/database"
	"github.com/waves2k/task-manager/internal/models"
)

type TodoRepository interface {
	Create(ctx context.Context, title string, completed bool) (*models.Todo, error)
	Update(ctx context.Context, params UpdateTodoParams) (*models.Todo, error)
	// GetByUserId(ctx context.Context, userId uuid.UUID) ([]models.Todo, error)
	GetById(ctx context.Context, todoId int) (*models.Todo, error)
	Delete(ctx context.Context, todoId int) error
}

type todoRepository struct {
	pool *pgxpool.Pool
}

func NewTodoService(pool *pgxpool.Pool) *todoRepository {
	return &todoRepository{
		pool: pool,
	}
}

type UpdateTodoParams struct {
	Id        int
	Title     *string
	Completed *bool
}

func makeSetQuery(params UpdateTodoParams) updateSetQuery {
	query := make([]string, 0)
	var result updateSetQuery
	result.ArgId = 1

	if params.Title != nil {
		query = append(query, fmt.Sprintf("title=$%d", result.ArgId))
		result.Args = append(result.Args, *params.Title)
		result.ArgId++
	}

	if params.Completed != nil {
		query = append(query, fmt.Sprintf("completed=$%d", result.ArgId))
		result.Args = append(result.Args, *params.Completed)
		result.ArgId++
	}

	query = append(query, "updated_at=CURRENT_TIMESTAMP")
	result.Sql = strings.Join(query, ", ")
	result.Args = append(result.Args, params.Id)

	return result
}

type updateSetQuery struct {
	Sql   string
	Args  []interface{}
	ArgId int
}

func (r *todoRepository) Update(ctx context.Context, params UpdateTodoParams) (*models.Todo, error) {
	intCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	setQuery := makeSetQuery(params)

	query := fmt.Sprintf(`
		UPDATE %s
		SET %s
		WHERE id=$%d
		RETURNING id, title, completed, created_at, updated_at
	`, database.TodoTable, setQuery.Sql, setQuery.ArgId)

	var todo models.Todo
	if err := r.pool.QueryRow(intCtx, query, setQuery.Args...).Scan(
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

func (r *todoRepository) Create(ctx context.Context, title string, completed bool) (*models.Todo, error) {
	intCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := fmt.Sprintf(`
		INSERT INTO %s (title, completed)
		VALUES ($1, $2)
		RETURNING id, title, completed, created_at, updated_at
	`, database.TodoTable)

	var todo models.Todo
	if err := r.pool.QueryRow(intCtx, query, title, completed).Scan(
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

func (r *todoRepository) GetById(ctx context.Context, todoId int) (*models.Todo, error) {
	intCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf(`
		SELECT id, title, completed, created_at, updated_at
		FROM %s
		WHERE id = $1
	`, database.TodoTable)

	var todo models.Todo
	if err := r.pool.QueryRow(intCtx, query, todoId).Scan(
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

func (r *todoRepository) Delete(ctx context.Context, todoId int) error {
	intCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE id = $1
	`, database.TodoTable)

	row, err := r.pool.Exec(intCtx, query, todoId)
	if err != nil {
		return err
	}
	if row.RowsAffected() == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// MAKE ORDER BY DESC
// func (r *todoRepository) GetAllByUserId(ctx context.Context, userId uuid.UUID) ([]models.Todo, error) {
// 	return nil, nil
// }
