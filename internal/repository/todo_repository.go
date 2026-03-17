package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/waves2k/task-manager/internal/database"
	apperrors "github.com/waves2k/task-manager/internal/errors"
	"github.com/waves2k/task-manager/internal/models"
)

type TodoRepository interface {
	Create(ctx context.Context, title string, completed bool, listId uuid.UUID) (*models.Todo, error)
	Update(ctx context.Context, params UpdateTodoParams) (*models.Todo, error)
	GetAll(ctx context.Context, listId, userId uuid.UUID) ([]models.Todo, error)
	GetById(ctx context.Context, todoId uuid.UUID) (*models.Todo, error)
	Delete(ctx context.Context, todoId uuid.UUID) error
}

type todoRepository struct {
	pool *pgxpool.Pool
}

func NewTodoService(pool *pgxpool.Pool) *todoRepository {
	return &todoRepository{
		pool: pool,
	}
}

func (r *todoRepository) Update(ctx context.Context, params UpdateTodoParams) (*models.Todo, error) {
	intCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	setQuery := makeTodoSetQuery(params)

	query := fmt.Sprintf(`
		UPDATE %s
		SET %s
		WHERE id=$%d
		RETURNING id, title, completed, created_at, updated_at
	`, database.TodosTable, setQuery.Sql, setQuery.ArgId)

	var todo models.Todo
	err := r.pool.QueryRow(intCtx, query, setQuery.Args...).Scan(
		&todo.Id,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && (pgErr.Code == "23505") {
		return nil, apperrors.NewAlreadyExistsError(
			apperrors.TodoAlreadyExistsMessage,
			err,
		)
	}
	if err != nil {
		return nil, apperrors.NewInternalError(
			apperrors.InternalErrorMessage,
			err,
		)
	}
	return &todo, nil
}

// USER_ID
func (r *todoRepository) Create(ctx context.Context, title string, completed bool, listId uuid.UUID) (*models.Todo, error) {
	intCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tx, err := r.pool.Begin(intCtx)
	if err != nil {
		return nil, apperrors.NewInternalError(
			apperrors.InternalErrorMessage,
			err,
		)
	}
	defer tx.Rollback(intCtx)

	mainQuery := fmt.Sprintf(`
		INSERT INTO %s (title, completed)
		VALUES ($1, $2)
		RETURNING id, title, completed, created_at, updated_at
	`, database.TodosTable)
	secondaryQuery := fmt.Sprintf(`
		INSERT INTO %s (list_id, todo_id)
		VALUES ($1, $2)
	`, database.ListsToTodosTable)

	var todo models.Todo
	if err := r.pool.QueryRow(intCtx, mainQuery, title, completed).Scan(
		&todo.Id,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && (pgErr.Code == "23505") {
			return nil, apperrors.NewAlreadyExistsError(
				apperrors.TodoNotFoundMessage,
				err,
			)
		}
		return nil, apperrors.NewInternalError(
			apperrors.InternalErrorMessage,
			err,
		)
	}

	_, err = tx.Exec(intCtx, secondaryQuery, listId, todo.Id)
	if err != nil {
		return nil, apperrors.NewInternalError(
			apperrors.InternalErrorMessage,
			err,
		)
	}
	if err = tx.Commit(intCtx); err != nil {
		return nil, apperrors.NewInternalError(
			apperrors.InternalErrorMessage,
			err,
		)
	}
	return &todo, nil
}

func (r *todoRepository) GetById(ctx context.Context, todoId uuid.UUID) (*models.Todo, error) {
	intCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf(`
		SELECT id, title, completed, created_at, updated_at
		FROM %s
		WHERE id = $1
	`, database.TodosTable)

	var todo models.Todo
	err := r.pool.QueryRow(intCtx, query, todoId).Scan(
		&todo.Id,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperrors.NewNotFoundError(
			apperrors.TodoNotFoundMessage,
			err,
		)
	}
	if err != nil {
		return nil, apperrors.NewInternalError(
			apperrors.InternalErrorMessage,
			err,
		)
	}
	return &todo, nil
}

func (r *todoRepository) Delete(ctx context.Context, todoId uuid.UUID) error {
	intCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE id = $1
	`, database.TodosTable)

	_, err := r.pool.Exec(intCtx, query, todoId)
	if errors.Is(err, pgx.ErrNoRows) {
		return apperrors.NewNotFoundError(
			apperrors.TodoNotFoundMessage,
			err,
		)
	}
	if err != nil {
		return apperrors.NewInternalError(
			apperrors.InternalErrorMessage,
			err,
		)
	}
	return nil
}

func (r *todoRepository) GetAll(ctx context.Context, listId, userId uuid.UUID) ([]models.Todo, error) {
	query := fmt.Sprintf(`
		SELECT td.id, td.title, td.completed, td.created_at, td.updated_at
		FROM %s td
		INNER JOIN %s li_td ON li_td.todo_id=td.id
		INNER JOIN %s us_li ON us_li.list_id = li_td.list_id
		WHERE li_td.list_id = $1 AND 
			us_li.user_id = $2
		ORDER BY td.created_at DESC
	`, database.TodosTable, database.ListsToTodosTable, database.UsersToListsTable)

	rows, err := r.pool.Query(ctx, query, listId, userId)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperrors.NewNotFoundError(
			apperrors.TodoNotFoundMessage,
			err,
		)
	}
	if err != nil {
		return nil, apperrors.NewInternalError(
			apperrors.InternalErrorMessage,
			err,
		)
	}
	todos := make([]models.Todo, 0)
	for rows.Next() {
		var todo models.Todo
		if err := rows.Scan(&todo.Id, &todo.Title, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
			return nil, apperrors.NewInternalError(
				apperrors.InternalErrorMessage,
				err,
			)
		}
		todos = append(todos, todo)
	}
	return todos, nil
}
