package service_test

import (
	"context"
	"errors"
	"testing"
	"user-service/mocks"
	"user-service/model"
	"user-service/service"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserService_CreateUser(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	svc := service.NewUserService(mockRepo)

	tests := []struct {
		name      string
		input     string
		mockFn    func()
		wantUser  model.User
		wantError bool
	}{
		{
			name:  "success",
			input: "Alice",
			mockFn: func() {
				mockRepo.EXPECT().
					CreateUser(ctx, "Alice").
					Return(model.User{ID: 1, Name: "Alice"}, nil)
			},
			wantUser:  model.User{ID: 1, Name: "Alice"},
			wantError: false,
		},
		{
			name:  "repo error",
			input: "Bob",
			mockFn: func() {
				mockRepo.EXPECT().
					CreateUser(ctx, "Bob").
					Return(model.User{}, errors.New("db error"))
			},
			wantUser:  model.User{},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			user, err := svc.CreateUser(ctx, tt.input)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantUser, user)
		})
	}
}

func TestUserService_GetUser(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	svc := service.NewUserService(mockRepo)

	tests := []struct {
		name      string
		input     uint64
		mockFn    func()
		wantUser  model.User
		wantError bool
	}{
		{
			name:  "success",
			input: 1,
			mockFn: func() {
				mockRepo.EXPECT().
					GetUser(ctx, uint64(1)).
					Return(model.User{ID: 1, Name: "Alice"}, nil)
			},
			wantUser:  model.User{ID: 1, Name: "Alice"},
			wantError: false,
		},
		{
			name:  "not found",
			input: 999,
			mockFn: func() {
				mockRepo.EXPECT().
					GetUser(ctx, uint64(999)).
					Return(model.User{}, errors.New("not found"))
			},
			wantUser:  model.User{},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			user, err := svc.GetUser(ctx, tt.input)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantUser, user)
		})
	}
}

func TestUserService_GetAllUsers(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	svc := service.NewUserService(mockRepo)

	tests := []struct {
		name      string
		page      int
		size      int
		mockFn    func()
		wantUsers []model.User
		wantError bool
	}{
		{
			name: "success",
			page: 2,
			size: 3,
			mockFn: func() {
				mockRepo.EXPECT().
					GetAllUsers(ctx, 3, 3).
					Return([]model.User{
						{ID: 4, Name: "D"},
						{ID: 5, Name: "E"},
					}, nil)
			},
			wantUsers: []model.User{
				{ID: 4, Name: "D"},
				{ID: 5, Name: "E"},
			},
			wantError: false,
		},
		{
			name: "repo error",
			page: 1,
			size: 2,
			mockFn: func() {
				mockRepo.EXPECT().
					GetAllUsers(ctx, 0, 2).
					Return(nil, errors.New("repo failure"))
			},
			wantUsers: nil,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			users, err := svc.GetAllUsers(ctx, tt.page, tt.size)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantUsers, users)
		})
	}
}

func TestUserService_GetUsersByIDs(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	svc := service.NewUserService(mockRepo)

	tests := []struct {
		name      string
		inputIDs  []uint64
		mockFn    func()
		wantUsers []model.User
		wantErr   bool
	}{
		{
			name:     "success - multiple users found",
			inputIDs: []uint64{1, 2},
			mockFn: func() {
				mockRepo.EXPECT().
					GetUserByIDs(ctx, []uint64{1, 2}).
					Return([]model.User{
						{ID: 1, Name: "Alice"},
						{ID: 2, Name: "Bob"},
					}, nil)
			},
			wantUsers: []model.User{
				{ID: 1, Name: "Alice"},
				{ID: 2, Name: "Bob"},
			},
			wantErr: false,
		},
		{
			name:      "success - empty input",
			inputIDs:  []uint64{},
			mockFn:    func() {}, // no repo call expected
			wantUsers: []model.User{},
			wantErr:   false,
		},
		{
			name:     "error from repo",
			inputIDs: []uint64{999},
			mockFn: func() {
				mockRepo.EXPECT().
					GetUserByIDs(ctx, []uint64{999}).
					Return(nil, errors.New("db error"))
			},
			wantUsers: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()

			users, err := svc.GetUsersByIDs(ctx, tt.inputIDs)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantUsers, users)
			}
		})
	}
}
