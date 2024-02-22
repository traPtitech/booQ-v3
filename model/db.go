package model

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

var allTables = []interface{}{
	Item{},
	Book{},
	Equipment{},
	Comment{},
	Transaction{},
	TransactionEquipment{},
	Tag{},
	Ownership{},
	Like{},
	File{},
}

type GormModel struct {
	ID        int       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

type GormModelWithoutID struct {
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

// EstablishConnection DBに接続する
func EstablishConnection() (error) {
	user := os.Getenv("MYSQL_USERNAME")
	if user == "" {
		user = "root"
	}

	pass := os.Getenv("MYSQL_PASSWORD")
	if pass == "" {
		pass = ""
	}

	host := os.Getenv("MYSQL_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("MYSQL_PORT")
	if port == "" {
		port = "3306"
	}

	dbname := os.Getenv("MYSQL_DATABASE")
	if dbname == "" {
		dbname = "booq-v3"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, dbname)+"?parseTime=true&loc=Asia%2FTokyo&charset=utf8mb4"
	_db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db = _db
	return err
}

func SetDBLoggerInfo() {
	db.Logger = db.Logger.LogMode(logger.Info)
}

// Migrate DBのマイグレーション
func Migrate() error {
	if err := db.AutoMigrate(allTables...); err != nil {
		return err
	}

	return nil
}