package sql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"golang.org/x/exp/slog"
	_ "gorm.io/driver/postgres"
)

var PagerDB *sql.DB
var PagerOrm *gorm.DB

func (dbConfig DatabaseConfigType) InitDatabase() (*sql.DB, *gorm.DB, error) {
	conString := getConnectionString(dbConfig.UserName, dbConfig.Password, dbConfig.Protocol, dbConfig.Host, dbConfig.Port, dbConfig.Database)
	db, err := gorm.Open("postgres", conString)
	if err != nil {
		slog.Error("DatabaseConnString", "conn_string", conString, "error", err)
		panic(err)
	}

	sqlDB := db.DB()

	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetConnMaxLifetime(time.Minute * 2)

	PagerDB = sqlDB
	PagerOrm = db

	return sqlDB, db, nil
}

func getConnectionString(username, password, protocol, host, port, dbname string) string {
	var conString string
	conString = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, dbname)
	//conString = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", "postgres", "password", "localhost", "5432", "pager_engine")

	return conString
}

func GetOrmQuearyable(ctx context.Context, tx interface{}) *gorm.DB {
	if tx == nil {
		return PagerOrm
	}
	return (tx).(*gorm.DB)
}

func TxRollBack(tx *gorm.DB, function string) {
	if tx == nil {
		return
	}
	if rollbackErr := tx.Rollback().Error; rollbackErr != nil {
		slog.Error("Error while rollback", "function", function, "error", rollbackErr)
	}
}
