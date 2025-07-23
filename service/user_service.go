package service

import (
	"context"
	"user-service/model"
	"user-service/repository"
)

// UserService defines user business logic contract.
//
//go:generate mockgen -source=user_service.go -destination=../mocks/mock_user_service.go -package=mocks
type UserService interface {
	CreateUser(ctx context.Context, name string) (model.User, error)
	GetUser(ctx context.Context, id uint64) (model.User, error)
	GetAllUsers(ctx context.Context, page, size int) ([]model.User, error)
	GetUsersByIDs(ctx context.Context, ids []uint64) ([]model.User, error)
}

// userServiceImpl is the actual implementation of UserService.
type userServiceImpl struct {
	repo repository.UserRepository
}

// NewUserService returns a UserService using the given UserRepository.
func NewUserService(repo repository.UserRepository) UserService {
	return &userServiceImpl{repo: repo}
}

func (s *userServiceImpl) CreateUser(ctx context.Context, name string) (model.User, error) {
	return s.repo.CreateUser(ctx, name)
}

func (s *userServiceImpl) GetUser(ctx context.Context, id uint64) (model.User, error) {
	return s.repo.GetUser(ctx, id)
}

func (s *userServiceImpl) GetAllUsers(ctx context.Context, page, size int) ([]model.User, error) {
	offset := (page - 1) * size
	return s.repo.GetAllUsers(ctx, offset, size)
}

// GetUsersByIDs fetches multiple users by their IDs
func (s *userServiceImpl) GetUsersByIDs(ctx context.Context, ids []uint64) ([]model.User, error) {
	if len(ids) == 0 {
		return make([]model.User, 0), nil
	}

	users, err := s.repo.GetUserByIDs(ctx, ids)

	if err != nil {
		return nil, err
	}

	return users, nil
}
