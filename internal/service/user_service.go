package service

import (
	"context"

	_ "github.com/golang-jwt/jwt"
	"github.com/waves2k/task-manager/internal/errors"
	"github.com/waves2k/task-manager/internal/models"
	"github.com/waves2k/task-manager/internal/repository"
)

type UserService interface {
	Register(ctx context.Context, email, password string) (*models.User, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type userService struct {
	userRepo       repository.UserRepository
	passwordHasher PasswordHasher
	tokenService   TokenService
}

func NewUserService(
	userRepo repository.UserRepository,
	passpasswordHasher PasswordHasher,
	tokenService TokenService) *userService {
	return &userService{
		userRepo:       userRepo,
		passwordHasher: passpasswordHasher,
		tokenService:   tokenService,
	}
}

func (s *userService) Register(ctx context.Context, email, password string) (*models.User, error) {
	hashedPassword := s.passwordHasher.GeneratePasswordHash(password)
	return s.userRepo.Create(ctx, email, hashedPassword)
}

func (s *userService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if !s.passwordHasher.VerifyPassword(password, user.PasswordHash) {
		return "", errors.NewWrongInputDataError(errors.WrongUserInputData, err)
	}
	token, err := s.tokenService.GenerateToken(user.Id)
	if err != nil {
		return "", err
	}
	return token, nil

}
