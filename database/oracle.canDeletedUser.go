package database

import (
	"context"
	"fmt"
	"time"

	"github.com/scch94/agentsDeleted/config"
	modeldb "github.com/scch94/agentsDeleted/models/db"
	"github.com/scch94/ins_log"
)

const (
	oracleIsUserParent = "SELECT * FROM user_Adm WHERE (user_parent_oid = :1 OR (VIEW_USER_PARENT_OID = :2 AND VIEW_USER_PARENT_OID <> OID)) AND tenant_oid = :3 AND ROWNUM = 1"
)

func isUserParentOracle(ctx context.Context, user *modeldb.UsersDb) error {
	ins_log.Tracef(ctx, "starting to check if the user %v can be deleted.", user.UserId)

	// Iniciamos el conteo en la base
	startTime := time.Now()

	ins_log.Tracef(ctx, "this is the QUERY: %s and the params: user_oid=%s, and tenant_oid=%v", oracleIsUserParent, user.UserOid, config.Config.Tenant)

	db := GetDb()

	rows, err := db.QueryContext(ctx, oracleIsUserParent, user.UserOid, user.UserOid, config.Config.Tenant)
	if err != nil {
		ins_log.Errorf(ctx, "query error %v", err)
		return err
	}
	defer rows.Close()

	if rows.Next() {
		ins_log.Infof(ctx, "User %v is a parent and cannot be deleted", user.UserId)
		user.CanDelete.AgentcanDeleted = false
		user.CanDelete.Reason = fmt.Sprintf("User %v is a parent and cannot be deleted", user.UserId)
	} else {
		ins_log.Infof(ctx, "User %v is not a parent and can be deleted", user.UserId)
		user.CanDelete.AgentcanDeleted = true
		user.CanDelete.Reason = fmt.Sprintf("User %v is not a parent and can be deleted", user.UserId)
	}
	// Verificar si hubo errores en el procesamiento de las filas
	if err = rows.Err(); err != nil {
		ins_log.Errorf(ctx, "error during row iteration: %v", err)
		return err
	}
	duration := time.Since(startTime)
	ins_log.Infof(ctx, "the query in the database took: %v", duration)

	return nil

}
