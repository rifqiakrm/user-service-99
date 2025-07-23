package model

// CreateUserRequest is a struct for CreateUser parameters
type CreateUserRequest struct {
	Name string `json:"name" binding:"required"`
}

// BatchFetchUsersRequest is the request payload for batch fetching users
type BatchFetchUsersRequest struct {
	UserIDs []uint64 `json:"user_ids" binding:"required"`
}

// BatchFetchUsersResponse is the response containing user data
type BatchFetchUsersResponse struct {
	Users []User `json:"users"`
}
