package repository_test

import (
	"context"
	"testing"
	"user-service/model"
	"user-service/repository"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return db
}

func TestUserRepo_CreateUser(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	repo := repository.NewUserRepo(db)

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid user", "John Doe", false},
		{"empty name", "", false}, // no validation on name
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.CreateUser(ctx, tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.input, user.Name)
				assert.NotZero(t, user.ID)
			}
		})
	}
}

func TestUserRepo_GetUser(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	repo := repository.NewUserRepo(db)

	createdUser, _ := repo.CreateUser(ctx, "Alice")

	tests := []struct {
		name     string
		userID   uint64
		wantErr  bool
		wantName string
	}{
		{"user exists", createdUser.ID, false, "Alice"},
		{"user not found", 9999, true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.GetUser(ctx, tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantName, user.Name)
			}
		})
	}
}

func TestUserRepo_GetUserByIDs(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	repo := repository.NewUserRepo(db)

	// Seed users
	user1, _ := repo.CreateUser(ctx, "Alice")
	user2, _ := repo.CreateUser(ctx, "Bob")

	tests := []struct {
		name      string
		ids       []uint64
		wantErr   bool
		wantCount int
		wantNames []string
	}{
		{
			name:      "users exist",
			ids:       []uint64{user1.ID, user2.ID},
			wantErr:   false,
			wantCount: 2,
			wantNames: []string{"Alice", "Bob"},
		},
		{
			name:      "some users exist",
			ids:       []uint64{user1.ID, 9999},
			wantErr:   false,
			wantCount: 1,
			wantNames: []string{"Alice"},
		},
		{
			name:      "no users exist",
			ids:       []uint64{9998, 9999},
			wantErr:   false,
			wantCount: 0,
			wantNames: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users, err := repo.GetUserByIDs(ctx, tt.ids)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Len(t, users, tt.wantCount)

			gotNames := make([]string, len(users))
			for i, u := range users {
				gotNames[i] = u.Name
			}
			assert.ElementsMatch(t, tt.wantNames, gotNames)
		})
	}
}

func TestUserRepo_GetAllUsers(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	repo := repository.NewUserRepo(db)

	// Seed 5 users
	for i := 1; i <= 5; i++ {
		repo.CreateUser(ctx, "User"+string(rune('A'+i-1)))
	}

	tests := []struct {
		name        string
		offset      int
		limit       int
		expectedLen int
	}{
		{"get first 3", 0, 3, 3},
		{"get next 2", 3, 2, 2},
		{"offset too high", 10, 5, 0},
		{"zero limit", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users, err := repo.GetAllUsers(ctx, tt.offset, tt.limit)
			assert.NoError(t, err)
			assert.Len(t, users, tt.expectedLen)
		})
	}
}
