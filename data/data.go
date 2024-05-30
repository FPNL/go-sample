package data

import (
	"fmt"

	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
	"go-micro.dev/v4/logger"

	"oltp/conf"
)

var ProviderSet = wire.NewSet(
	NewData,
	NewSQLX,

	NewGreeterRepo,
)

// Data 存取資料，包含 mysql, redis...
type Data struct {
	db *sqlx.DB
	// cache redis.UniversalClient
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
func NewData(db *sqlx.DB) (*Data, func(), error) {
	cleanup := func() {
		_ = db.Close()
		logger.Info("closing the data resources")
	}
	return &Data{
		db: db,
	}, cleanup, nil
}

func source(conf *conf.Data) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		conf.Mysql.User, conf.Mysql.Pwd, conf.Mysql.URL,
		conf.Mysql.Port, conf.Mysql.DbName)
}
