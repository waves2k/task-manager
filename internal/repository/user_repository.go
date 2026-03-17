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

type UserRepository interface {
	Create(ctx context.Context, email string, passwordHash string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetById(ctx context.Context, id uuid.UUID) (*models.User, error)
	IsExists(ctx context.Context, id uuid.UUID) bool
}

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *userRepository {
	return &userRepository{
		pool: pool,
	}
}

func (r *userRepository) Create(ctx context.Context, email string, passwordHash string) (*models.User, error) {
	intCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := fmt.Sprintf(`
		INSERT INTO %s (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, email, password_hash, created_at, updated_at
	`, database.UsersTable)

	var user models.User
	if err := r.pool.QueryRow(intCtx, query, email, passwordHash).Scan(
		&user.Id,
		&user.Email,
		&user.PasswordHash, &user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && (pgErr.Code == "23505") {
			return nil, apperrors.NewAlreadyExistsError(
				apperrors.UserAlreadyExistsMessage,
				err,
			)
		}
		return nil, apperrors.NewInternalError(
			apperrors.InternalErrorMessage,
			err,
		)
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := fmt.Sprintf(`
		SELECT id, email, password_hash, created_at, updated_at
		FROM %s
		WHERE email=$1
	`, database.UsersTable)

	var user models.User
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.Id,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperrors.NewNotFoundError(
			apperrors.UserNotFoundMessage,
			err,
		)
	}
	if err != nil {
		return nil, apperrors.NewInternalError(
			apperrors.InternalErrorMessage,
			err,
		)
	}
	return &user, nil
}

func (r *userRepository) GetById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := fmt.Sprintf(`
		SELECT id, email, password_hash, created_at, updated_at
		FROM %s
		WHERE id=$1
	`, database.UsersTable)

	var user models.User
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.Id,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperrors.NewNotFoundError(
			apperrors.UserNotFoundMessage,
			err,
		)
	}
	if err != nil {
		return nil, apperrors.NewInternalError(
			apperrors.InternalErrorMessage,
			err,
		)
	}
	return &user, nil
}

func (r *userRepository) IsExists(ctx context.Context, id uuid.UUID) bool {
	query := fmt.Sprintf(`
		SELECT id
		FROM %s
		WHERE id=$1
	`, database.UsersTable)

	if err := r.pool.QueryRow(ctx, query, id).Scan(); err != nil {
		return false
	}
	return true
}
