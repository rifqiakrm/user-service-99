package model

// User represents a user in the system.
type User struct {
	ID        uint64 `json:"id" gorm:"primaryKey"` // Unique user ID
	Name      string `json:"name"`                 // Full name of the user
	CreatedAt int64  `json:"created_at"`           // Timestamp in microseconds
	UpdatedAt int64  `json:"updated_at"`           // Timestamp in microseconds
}
