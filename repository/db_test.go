package repository

import (
	"context"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mariadb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupTestDB(ctx context.Context, t *testing.T) *gorm.DB {
	container, err := mariadb.Run(ctx,
		"mariadb:12.0.2",
		mariadb.WithDatabase("testdb"),
		mariadb.WithUsername("testuser"),
		mariadb.WithPassword("testpass"),
	)
	if err != nil {
		t.Fatalf("Failed to start MariaDB container: %v", err)
	}

	t.Cleanup(func() {
		if err := testcontainers.TerminateContainer(container); err != nil {
			t.Fatalf("Failed to terminate container: %v", err)
		}
	})

	conn, err := container.ConnectionString(ctx,
		"charset=utf8mb4",
		"parseTime=true",
		"loc=Asia%2FTokyo",
		"allowNativePasswords=true",
	)
	if err != nil {
		t.Fatalf("Failed to get connection string: %v", err)
	}

	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	if err := Migrate(db); err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}
