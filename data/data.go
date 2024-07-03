package data

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
	"go-micro.dev/v4/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/fpnl/go-sample/conf"
)

var ProviderSet = wire.NewSet(
	NewData,
	NewSQLX,
	NewGORM,

	NewGreeterRepo,
)

// Data 存取資料，包含 mysql, redis...
type Data struct {
	logger *slog.Logger
	db     *gorm.DB
	// cache redis.UniversalClient
}

func NewGORM(conf *conf.Data) (*gorm.DB, error) {
	sqlDB, err := sql.Open("mysql", source(conf)+"?parseTime=true&loc=Local")
	if err != nil {
		return nil, fmt.Errorf("sql.Open() : %v", err)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("gorm.Open() : %v", err)
	}

	return gormDB, nil
}

func NewSQLX(conf *conf.Data) *sqlx.DB {
	db, err := sqlx.Open("mysql", source(conf)+"?parseTime=true&loc=Local")
	if err != nil {
		logger.Fatalf("sqlx.Open() : %v", err)
	}

	err = db.Ping()
	if err != nil {
		logger.Fatalf("db.Ping() : %v", err)
	}

	db.SetMaxIdleConns(conf.Mysql.Conn)
	db.SetMaxOpenConns(conf.Mysql.Conn)

	return db
}

// NewData .
func NewData(db *gorm.DB, logger *slog.Logger) (*Data, func(), error) {
	cleanup := func() {
		db, _ := db.DB()
		db.Close()
		logger.Info("closing the data resources")
	}
	return &Data{
		db:     db,
		logger: logger,
	}, cleanup, nil
}

func source(conf *conf.Data) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		conf.Mysql.User, conf.Mysql.Pwd, conf.Mysql.URL,
		conf.Mysql.Port, conf.Mysql.DbName)
}
