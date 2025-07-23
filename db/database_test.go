// db/database_test.go
package db_test

import (
	"os"
	"testing"

	"user-service/db"
	"user-service/model"

	"github.com/stretchr/testify/assert"
)

func TestInitDB(t *testing.T) {
	tests := []struct {
		name      string
		dbPath    string
		expectErr bool
	}{
		{
			name:      "successfully connects and migrates",
			dbPath:    "test_user.db",
			expectErr: false,
		},
		{
			name:      "fail to connect to invalid db path",
			dbPath:    "/invalid_path/test_user.db",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Remove(tt.dbPath)

			dbInstance, err := db.InitDB(tt.dbPath)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, dbInstance)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, dbInstance)
				assert.True(t, dbInstance.Migrator().HasTable(&model.User{}))
			}
		})
	}
}
