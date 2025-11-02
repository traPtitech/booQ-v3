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

type DB struct {
	db *gorm.DB
}

var allTables = []interface{}{
	item{},
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

func NewDB(db *gorm.DB) *DB {
	return &DB{db: db}
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

func (d *DB) SetLoggerInfo() {
	d.db.Logger = d.db.Logger.LogMode(logger.Info)
}

func (d *DB) Migrate() error {
	if err := d.db.AutoMigrate(allTables...); err != nil {
		return err
	}
	return nil
}
