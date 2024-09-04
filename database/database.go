package database

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/godror/godror" // Importa el driver para Oracle
	_ "github.com/lib/pq"        // Importa el driver para PostgreSQL
	"github.com/scch94/agentsDeleted/config"
	modeldb "github.com/scch94/agentsDeleted/models/db"
	modelUtils "github.com/scch94/agentsDeleted/models/utils"
	"github.com/scch94/ins_log"
)

var (
	DB             *sql.DB
	dbOnce         sync.Once
	databaseEngine string
)

func InitDb(ctx context.Context) error {
	ctx = ins_log.SetPackageNameInContext(ctx, "databaseConnection")
	if config.Config.DatabaseEngine == DBPOSTGRES {
		err := NewPostgresDb(ctx, config.Config.DatabaseConnectionString)
		if err != nil {
			return err
		}
		return nil
	} else if config.Config.DatabaseEngine == DBORACLE {
		err := NewOracleDb(ctx, config.Config.DatabaseConnectionString)
		if err != nil {
			return err
		}
		return nil
	}
	ins_log.Errorf(ctx, "please check the configuration file in databaseengine the values accepted are postgresql or oraclesql")
	return fmt.Errorf("please check the configuration file in databaseengine the values accepted are postgresql or oraclesql")

}
func NewPostgresDb(ctx context.Context, connectionString string) error {
	//"postgres://xxepin:migracion@digicel-dev-flex.postgres.database.azure.com:5432/xxepin?sslmode=require",
	databaseEngine = config.Config.DatabaseEngine
	var initErr error
	dbOnce.Do(func() {
		var err error

		DB, err = sql.Open("postgres", connectionString)
		if err != nil {
			ins_log.Fatalf(ctx, "cant open postgres database with string connection %v and the error is: %v", connectionString, err)
			initErr = err
			return
		}
		DB.SetConnMaxIdleTime(1800)
		DB.SetConnMaxLifetime(3600)
		DB.SetMaxOpenConns(1000)
		DB.SetMaxIdleConns(1000)

		if err = DB.Ping(); err != nil {
			ins_log.Fatalf(ctx, "cant do ping to  postgres database error : %v", err)
			initErr = err
			return
		}
		ins_log.Info(ctx, "connected to oracle database")
	})
	return initErr
}
func NewOracleDb(ctx context.Context, connectionString string) error {
	databaseEngine = config.Config.DatabaseEngine
	//coneection := "EPIN_NEW/epin@192.168.0.157:1521/EPIN"
	var initErr error
	dbOnce.Do(func() {
		var err error

		DB, err = sql.Open("godror", connectionString)
		if err != nil {
			ins_log.Fatalf(ctx, "cant open postgres database with string connection %v and the error is: %v", connectionString, err)
			initErr = err
			return
		}
		DB.SetConnMaxIdleTime(1800)
		DB.SetConnMaxLifetime(3600)
		DB.SetMaxOpenConns(1000)
		DB.SetMaxIdleConns(1000)

		if err = DB.Ping(); err != nil {
			ins_log.Fatalf(ctx, "cant do ping to the oracle database error : %v", err)
			initErr = err
			return
		}
		ins_log.Info(ctx, "connected to oracle database")
	})
	return initErr
}
func newDb(ctx context.Context, databaseEngine string) error {
	connectionString := config.Config.DatabaseConnectionString
	var initErr error
	dbOnce.Do(func() {
		var err error

		DB, err = sql.Open(databaseEngine, connectionString)
		if err != nil {
			ins_log.Fatalf(ctx, "cant open database with string connection %v and the error is: %v", connectionString, err)
			initErr = err
			return
		}
		DB.SetConnMaxIdleTime(1800)
		DB.SetConnMaxLifetime(3600)
		DB.SetMaxOpenConns(1000)
		DB.SetMaxIdleConns(1000)

		if err = DB.Ping(); err != nil {
			ins_log.Fatalf(ctx, "cant do ping to the database error : %v", err)
			initErr = err
			return
		}
		ins_log.Info(ctx, "connected to database")
	})
	return initErr
}

func GetDb() *sql.DB {
	return DB
}

// funciones que solo sirven de gateways para ver que motor estan usando y que consultas son.
func IsAgentParent(ctx context.Context, agent *modelUtils.Agents) error {
	if agent.AgentOid != "" {
		if databaseEngine == "oraclesql" {
			err := IsAgentParentOracle(ctx, agent)
			if err != nil {
				ins_log.Errorf(ctx, "error in IsAgentParentOracle() err: %v", err)
				return err
			}
		} else {
			err := IsAgentParentPostgres(ctx, agent)
			if err != nil {
				ins_log.Errorf(ctx, "error in IsAgentParentPostgres() err: %v", err)
				return err
			}
		}

	} else {
		ins_log.Infof(ctx, "the agent with id %v didnt exist", agent.AgentId)
		agent.CanDelete.AgentcanDeleted = false
		agent.CanDelete.Reason = fmt.Sprintf("the agent with id %v didnt exist", agent.AgentId)
	}

	return nil
}

func GetUsers(ctx context.Context, agent *modelUtils.Agents) ([]modeldb.UsersDb, error) {
	var users []modeldb.UsersDb
	var err error
	//CHEQUEAMOS SI EL AGENTE EXISTE
	if agent.AgentOid != "" {
		if databaseEngine == "oraclesql" {
			users, err = GetUsersOracle(ctx, agent)
			if err != nil {
				ins_log.Errorf(ctx, "error in IsAgentParentOracle() err: %v", err)
				return nil, err
			}
		} else {
			users, err = GetUsersPostgres(ctx, agent)
			if err != nil {
				ins_log.Errorf(ctx, "error in IsAgentParentPostgres() err: %v", err)
				return nil, err
			}
		}
	} else {
		ins_log.Infof(ctx, "the agent with id %v didnt exist", agent.AgentId)
		agent.CanDelete.AgentcanDeleted = false
		agent.CanDelete.Reason = fmt.Sprintf("the agent with id %v didnt exist", agent.AgentId)
	}
	return users, nil
}

func GetMsisdn(ctx context.Context, agent *modelUtils.Agents) ([]modeldb.MsisdnDb, error) {
	var agentsMsisdn []modeldb.MsisdnDb
	var err error
	ins_log.Tracef(ctx, "this is the database engine %s", databaseEngine)
	if databaseEngine == "oraclesql" {
		agentsMsisdn, err = GetMsisdnOracle(ctx, agent)
		if err != nil {
			ins_log.Errorf(ctx, "error in GetMsisdnOracle() err: %v", err)
			return nil, err
		}
	} else {
		agentsMsisdn, err = GetMsisdnPostgres(ctx, agent)
		if err != nil {
			ins_log.Errorf(ctx, "error in GetMsisdnPostgres() err: %v", err)
			return nil, err
		}
	}
	return agentsMsisdn, nil

}
