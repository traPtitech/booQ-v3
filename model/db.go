package model

import (
	"fmt"
	"os"
	"time"

	"github.com/go-testfixtures/testfixtures/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db       *gorm.DB
	fixtures *testfixtures.Loader
)

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
	ID        int       `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type GormModelWithoutID struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// EstablishConnection DBに接続する
func EstablishConnection() error {
	user := os.Getenv("MYSQL_USERNAME")
	if user == "" {
		user = "root"
	}

	pass := os.Getenv("MYSQL_PASSWORD")
	if pass == "" {
		pass = "password"
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

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, dbname) + "?parseTime=true&loc=Asia%2FTokyo&charset=utf8mb4"
	_db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return err
	}

	db = _db
	return nil
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

// テスト用DBの稼働。model, routerのテストで用いる
func SetUpTestDB() {
	err := EstablishConnection()
	if err != nil {
		panic(err)
	}

	err = Migrate()
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	dirFixtures := "../testdata/fixtures"
	fixtures, err = testfixtures.New(
		testfixtures.Database(sqlDB),        // You database connection
		testfixtures.Dialect("mysql"),       // Available: "postgresql", "timescaledb", "mysql", "mariadb", "sqlite" and "sqlserver"
		testfixtures.Directory(dirFixtures), // The directory containing the YAML files
	)

	if err != nil {
		wd, _ := os.Getwd()
		panic(fmt.Errorf("%v %v", err, wd))
	}
}

func PrepareTestDatabase() {
	if err := fixtures.Load(); err != nil {
		panic(err)
	}
}
