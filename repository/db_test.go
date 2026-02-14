package repository

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/labstack/gommon/log"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mariadb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	ctx := context.Background()
	container, err := mariadb.Run(ctx,
		"mariadb:12.0.2",
		mariadb.WithDatabase("testdb"),
		mariadb.WithUsername("testuser"),
		mariadb.WithPassword("testpass"),
	)
	if err != nil {
		log.Fatalf("Failed to start MariaDB container: %v", err)
	}

	conn, err := container.ConnectionString(ctx,
		"charset=utf8mb4",
		"parseTime=true",
		"loc=Asia%2FTokyo",
		"allowNativePasswords=true",
	)
	if err != nil {
		log.Fatalf("Failed to get connection string: %v", err)
	}

	db, err = gorm.Open(mysql.Open(conn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := Migrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	code := m.Run()

	if err := testcontainers.TerminateContainer(container); err != nil {
		log.Printf("Failed to terminate container: %v", err)
	}

	os.Exit(code)
}

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	t.Cleanup(func() {
		if err := db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
			t.Fatalf("Failed to disable foreign key checks: %v", err)
		}

		var tables []string
		if err := db.Raw("SHOW TABLES").Scan(&tables).Error; err != nil {
			t.Fatalf("Failed to get table names: %v", err)
		}

		for _, table := range tables {
			if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", table)).Error; err != nil {
				t.Fatalf("Failed to truncate table %s: %v", table, err)
			}
		}

		if err := db.Exec("SET FOREIGN_KEY_CHECKS = 1").Error; err != nil {
			t.Fatalf("Failed to enable foreign key checks: %v", err)
		}
	})

	return db
}
