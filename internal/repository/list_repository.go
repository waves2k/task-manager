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

type ListRepository interface {
	Create(ctx context.Context, title, description string, userId uuid.UUID) (*models.List, error)
	GetById(ctx context.Context, listId, userId uuid.UUID) (*models.List, error)
	GetAll(ctx context.Context, userId uuid.UUID) ([]models.List, error)
	Update(ctx context.Context, params UpdateListParams) (*models.List, error)
	Delete(ctx context.Context, listId, userId uuid.UUID) error
	IsExists(ctx context.Context, id uuid.UUID) bool
}

type listRepository struct {
	pool *pgxpool.Pool
}

func NewListRepository(pool *pgxpool.Pool) *listRepository {
	return &listRepository{
		pool: pool,
	}
}

func (r *listRepository) Create(ctx context.Context, title, description string, userId uuid.UUID) (*models.List, error) {
	intCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	listsQuery := fmt.Sprintf(`
		INSERT INTO %s (title, description)
		VALUES ($1, $2)
		RETURNING id, title, description, created_at, updated_at
	`, database.ListsTable)

	usersListsQuery := fmt.Sprintf(`
		INSERT INTO %s (user_id, list_id)
		VALUES ($1, $2)
	`, database.UsersToListsTable)

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, apperrors.NewInternalError(apperrors.InternalErrorMessage, err)
	}
	defer tx.Rollback(ctx)

	var list models.List
	if err := tx.QueryRow(intCtx, listsQuery, title, description).Scan(
		&list.Id,
		&list.Title,
		&list.Description,
		&list.CreatedAt,
		&list.UpdatedAt,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && (pgErr.Code == "23505") {
			return nil, apperrors.NewAlreadyExistsError(
				apperrors.ListAlreadyExistsMessage,
				err,
			)
		}
		return nil, apperrors.NewInternalError(
			apperrors.InternalErrorMessage,
			err,
		)
	}

	_, err = tx.Exec(intCtx, usersListsQuery, list.Id)
	if err != nil {
		return nil, apperrors.NewInternalError(apperrors.InternalErrorMessage, err)
	}

	tx.Commit(ctx)
	return &list, nil
}

func (r *listRepository) GetById(ctx context.Context, listId, userId uuid.UUID) (*models.List, error) {
	query := fmt.Sprintf(`
		SELECT li.id, li.title, li.description, li.created_at, li.updated_at
		FROM %s li
		INNER JOIN %s ul ON li.id = ul.list_id
		WHERE li.id = $1
		AND ul.user_id = $2
	`, database.ListsTable, database.UsersToListsTable)

	var list models.List
	err := r.pool.QueryRow(ctx, query, listId, userId).Scan(
		&list.Id,
		&list.Title,
		&list.Description,
		&list.CreatedAt,
		&list.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperrors.NewNotFoundError(
			apperrors.ListNotFoundMessage,
			err,
		)
	}
	if err != nil {
		return nil, apperrors.NewInternalError(
			apperrors.InternalErrorMessage,
			err,
		)
	}
	return &list, nil
}

func (r *listRepository) GetAll(ctx context.Context, userId uuid.UUID) ([]models.List, error) {
	query := fmt.Sprintf(`
		SELECT li.id, li.title, li.description, li.created_at, li.updated_at
		FROM %s li
		INNER JOIN %s ul ON li.id = ul.list_id
		WHERE ul.user_id = $1
		ORDER BY li.created_at DESC
	`, database.ListsTable, database.UsersToListsTable)

	rows, err := r.pool.Query(ctx, query, userId)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperrors.NewNotFoundError(
			apperrors.ListNotFoundMessage,
			err,
		)
	}
	if err != nil {
		return nil, apperrors.NewInternalError(
			apperrors.InternalErrorMessage,
			err,
		)
	}

	lists := make([]models.List, 0)
	for rows.Next() {
		var list models.List
		if err := rows.Scan(&list.Id, &list.Title, &list.Description, &list.CreatedAt, &list.UpdatedAt); err != nil {
			return nil, apperrors.NewInternalError(
				apperrors.InternalErrorMessage,
				err,
			)
		}
		lists = append(lists, list)
	}
	rows.Close()
	return lists, nil
}

func (r *listRepository) Update(ctx context.Context, params UpdateListParams) (*models.List, error) {
	setQuert := makeListSetQuery(params)
	query := fmt.Sprintf(`
		UPDATE %s li
		SET %s 
		FROM %s ul
		WHERE ul.list_id = li.id
		AND li.id = $%d
		AND ul.user_id = $%d
		RETURNING li.id, li.title, li.description, li.created_at, li.updated_at
	`, database.ListsTable, setQuert.Sql, setQuert.ArgId, setQuert.ArgId+1)

	var list models.List
	err := r.pool.QueryRow(ctx, query, setQuert.Args...).Scan(
		&list.Id,
		&list.Title,
		&list.Description,
		&list.CreatedAt,
		&list.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperrors.NewNotFoundError(
			apperrors.ListNotFoundMessage,
			err,
		)
	}
	if err != nil {
		return nil, apperrors.NewInternalError(
			apperrors.InternalErrorMessage,
			err,
		)
	}
	return &list, nil
}

func (r *listRepository) Delete(ctx context.Context, listId, userId uuid.UUID) error {
	intCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	mainQuery := fmt.Sprintf(`
		DELETE FROM %s
		WHERE id = $1
	`, database.ListsTable)

	secondaryQuery := fmt.Sprintf(`
		DELETE FROM %s td
		USING %s li_td 
		WHERE td.id = li_td.todo_id
		AND li.list_id = $1
	`, database.TodosTable, database.ListsToTodosTable)

	tx, err := r.pool.Begin(intCtx)
	if err != nil {
		return apperrors.NewInternalError(apperrors.InternalErrorMessage, err)
	}
	defer tx.Rollback(intCtx)

	_, err = tx.Exec(intCtx, mainQuery, listId)
	if errors.Is(err, pgx.ErrNoRows) {
		return apperrors.NewNotFoundError(
			apperrors.ListNotFoundMessage,
			err,
		)
	}
	if err != nil {
		return apperrors.NewInternalError(apperrors.InternalErrorMessage, err)
	}

	_, err = tx.Exec(intCtx, secondaryQuery, listId)
	if errors.Is(err, pgx.ErrNoRows) {
		return apperrors.NewNotFoundError(
			apperrors.TodoNotFoundMessage,
			err,
		)
	}
	if err != nil {
		return apperrors.NewInternalError(apperrors.InternalErrorMessage, err)
	}

	err = tx.Commit(intCtx)
	if err != nil {
		return apperrors.NewInternalError(apperrors.InternalErrorMessage, err)
	}
	return nil
}

func (r *listRepository) IsExists(ctx context.Context, id uuid.UUID) bool {
	query := fmt.Sprintf(`
		SELECT id
		FROM %s
		WHERE id=$1
	`, database.ListsTable)

	if err := r.pool.QueryRow(ctx, query, id).Scan(); err != nil {
		return false
	}
	return true
}
