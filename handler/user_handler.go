package handler

import (
	"net/http"
	"strconv"
	"user-service/model"
	"user-service/service"

	"github.com/gin-gonic/gin"
)

// UserHandler handles HTTP requests related to users.
type UserHandler struct {
	Svc service.UserService
}

// NewUserHandler initializes the user handler with service dependency.
func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{Svc: svc}
}

// CreateUser handles POST /users
// Accepts form data and creates a new user.
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req model.CreateUserRequest

	// Bind and validate input
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service layer
	user, err := h.Svc.CreateUser(c.Request.Context(), req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"result": true, "user": user})
}

// GetUser handles GET /users/:id
// Returns a user by ID.
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "error": "invalid id"})
		return
	}
	user, err := h.Svc.GetUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"result": false, "error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": true, "user": user})
}

// GetAllUsers handles GET /users
// Supports pagination via page_num & page_size query parameters.
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	pageNum, _ := strconv.Atoi(c.DefaultQuery("page_num", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	users, err := h.Svc.GetAllUsers(c.Request.Context(), pageNum, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"result": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": true, "users": users})
}

// BatchFetchUsers handles POST /users/batch to fetch multiple users by IDs
func (h *UserHandler) BatchFetchUsers(c *gin.Context) {
	var req model.BatchFetchUsersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	users, err := h.Svc.GetUsersByIDs(c.Request.Context(), req.UserIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": true, "users": users})
}
