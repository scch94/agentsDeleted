package database

import (
	"context"
	"database/sql"
	"sync"

	_ "github.com/lib/pq"
	"github.com/scch94/agentsDeleted/config"
	"github.com/scch94/ins_log"
)

var (
	DB     *sql.DB
	dbOnce sync.Once
)

func InitDb(ctx context.Context) {
	NewPostgresDb(ctx)
}
func NewPostgresDb(ctx context.Context) {
	dbOnce.Do(func() {
		ctx = ins_log.SetPackageNameInContext(ctx, "databaseConnection")
		var err error

		DB, err = sql.Open("postgres", config.Config.DatabaseConnectionString)
		if err != nil {
			ins_log.Fatalf(ctx, "cant open postgres database with string connection %v and the error is: %v", config.Config.DatabaseConnectionString, err)
		}
		DB.SetConnMaxIdleTime(1800)
		DB.SetConnMaxLifetime(3600)
		DB.SetMaxOpenConns(2)
		DB.SetMaxIdleConns(2)

		if err = DB.Ping(); err != nil {
			ins_log.Fatalf(ctx, "cant do ping to database error : %v", err)
		}
		ins_log.Info(ctx, "connected to postgres database")
	})

}

func GetDb() *sql.DB {
	return DB
}
