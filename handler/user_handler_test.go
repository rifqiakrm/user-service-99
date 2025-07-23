package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"user-service/mocks"
	"user-service/model"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func setupRouter(h *UserHandler) *gin.Engine {
	r := gin.Default()
	r.POST("/users", h.CreateUser)
	r.GET("/users/:id", h.GetUser)
	r.GET("/users", h.GetAllUsers)
	r.POST("/users/batch", h.BatchFetchUsers)
	return r
}

func TestCreateUser(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockUserService(ctrl)
	handler := NewUserHandler(mockSvc)

	tests := []struct {
		name           string
		body           gin.H
		mockFunc       func()
		expectedStatus int
	}{
		{
			name: "success",
			body: gin.H{"name": "Alice"},
			mockFunc: func() {
				mockSvc.EXPECT().
					CreateUser(ctx, "Alice").
					Return(model.User{ID: 1, Name: "Alice"}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "missing name",
			body: gin.H{},
			mockFunc: func() {
				// No call expected
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "internal error",
			body: gin.H{"name": "Bob"},
			mockFunc: func() {
				mockSvc.EXPECT().
					CreateUser(ctx, "Bob").
					Return(model.User{}, errors.New("create failed"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	router := setupRouter(handler)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			bodyBytes, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestGetUser(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockUserService(ctrl)
	handler := NewUserHandler(mockSvc)
	router := setupRouter(handler)

	tests := []struct {
		name           string
		paramID        string
		mockFunc       func()
		expectedStatus int
	}{
		{
			name:    "valid ID found",
			paramID: "1",
			mockFunc: func() {
				mockSvc.EXPECT().
					GetUser(ctx, uint64(1)).
					Return(model.User{ID: 1, Name: "Alice"}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:    "invalid ID param",
			paramID: "abc",
			mockFunc: func() {
				// Should not call mock
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:    "user not found",
			paramID: "10",
			mockFunc: func() {
				mockSvc.EXPECT().
					GetUser(ctx, uint64(10)).
					Return(model.User{}, errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			req, _ := http.NewRequest(http.MethodGet, "/users/"+tt.paramID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestGetAllUsers(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockUserService(ctrl)
	handler := NewUserHandler(mockSvc)
	router := setupRouter(handler)

	tests := []struct {
		name           string
		mockFunc       func()
		expectedStatus int
	}{
		{
			name: "success get all",
			mockFunc: func() {
				mockSvc.EXPECT().
					GetAllUsers(ctx, 1, 10).
					Return([]model.User{
						{ID: 1, Name: "Alice"},
						{ID: 2, Name: "Bob"},
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "internal server error",
			mockFunc: func() {
				mockSvc.EXPECT().
					GetAllUsers(ctx, 1, 10).
					Return(nil, errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			req, _ := http.NewRequest(http.MethodGet, "/users", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestBatchFetchUsers(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockUserService(ctrl)
	handler := NewUserHandler(mockSvc)
	router := setupRouter(handler)

	tests := []struct {
		name           string
		requestBody    string
		mockFunc       func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "success",
			requestBody: `{"user_ids": [1, 2]}`,
			mockFunc: func() {
				mockSvc.EXPECT().
					GetUsersByIDs(ctx, []uint64{1, 2}).
					Return([]model.User{
						{ID: 1, Name: "Alice"},
						{ID: 2, Name: "Bob"},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"result":true`, // partial match
		},
		{
			name:           "invalid request body",
			requestBody:    `{"user_ids": "not-an-array"}`,
			mockFunc:       func() {}, // no mock expected
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"error":"invalid request"`,
		},
		{
			name:        "internal server error",
			requestBody: `{"user_ids": [1, 2]}`,
			mockFunc: func() {
				mockSvc.EXPECT().
					GetUsersByIDs(gomock.Any(), []uint64{1, 2}).
					Return(nil, errors.New("something went wrong"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `"error":"something went wrong"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			req, _ := http.NewRequest(http.MethodPost, "/users/batch", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)
		})
	}
}
