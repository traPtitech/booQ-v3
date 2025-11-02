package repository

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/mysql"
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
	user := os.Getenv("MYSQL_USERNAME")
	pass := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	name := os.Getenv("MYSQL_DATABASE")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Asia%2FTokyo&charset=utf8mb4", user, pass, host, port, name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
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
