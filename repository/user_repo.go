package repository

import (
	"context"
	"time"
	"user-service/model"

	"gorm.io/gorm"
)

// UserRepository defines the contract for user data access layer.
//
//go:generate mockgen -source=user_repo.go -destination=../mocks/mock_user_repo.go -package=mocks
type UserRepository interface {
	CreateUser(ctx context.Context, name string) (model.User, error)
	GetUser(ctx context.Context, id uint64) (model.User, error)
	GetUserByIDs(ctx context.Context, ids []uint64) ([]model.User, error)
	GetAllUsers(ctx context.Context, offset, limit int) ([]model.User, error)
}

// userRepoImpl is the concrete implementation of UserRepository using GORM.
type userRepoImpl struct {
	DB *gorm.DB
}

// NewUserRepo initializes the db and returns a UserRepository.
func NewUserRepo(db *gorm.DB) UserRepository {
	return &userRepoImpl{DB: db}
}

func (r *userRepoImpl) CreateUser(ctx context.Context, name string) (model.User, error) {
	now := time.Now().UnixMicro()
	user := model.User{
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
	result := r.DB.WithContext(ctx).Create(&user)
	return user, result.Error
}

func (r *userRepoImpl) GetUser(ctx context.Context, id uint64) (model.User, error) {
	var user model.User
	result := r.DB.WithContext(ctx).First(&user, id)
	return user, result.Error
}

func (r *userRepoImpl) GetUserByIDs(ctx context.Context, ids []uint64) ([]model.User, error) {
	user := make([]model.User, 0)

	result := r.DB.WithContext(ctx).Where("id in (?)", ids).Find(&user)
	return user, result.Error
}

func (r *userRepoImpl) GetAllUsers(ctx context.Context, offset, limit int) ([]model.User, error) {
	var users []model.User
	result := r.DB.WithContext(ctx).Order("created_at desc").Offset(offset).Limit(limit).Find(&users)
	return users, result.Error
}
