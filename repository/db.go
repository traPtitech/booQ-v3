package repository

import (
	"fmt"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var allTables = []interface{}{
	item{},
	tag{},
	like{},
	book{},
	equipment{},
	file{},
	ownership{},
	transaction{},
}

type GormModel struct {
	ID        int `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type GormModelWithoutID struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

func EstablishConnection() (*gorm.DB, error) {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, err
	}

	cfg := mysql.Config{
		User:                 os.Getenv("MYSQL_USERNAME"),
		Passwd:               os.Getenv("MYSQL_PASSWORD"),
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_PORT")),
		DBName:               os.Getenv("MYSQL_DATABASE"),
		ParseTime:            true,
		Loc:                  loc,
		Collation:            "utf8mb4_general_ci",
		AllowNativePasswords: true,
	}

	db, err := gorm.Open(gormMysql.Open(cfg.FormatDSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(allTables...)
}

func SetLoggerInfo(db *gorm.DB) {
	db.Logger = db.Logger.LogMode(logger.Info)
}
