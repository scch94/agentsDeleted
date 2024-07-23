package database

import (
	"context"
	"time"

	"github.com/scch94/agentsDeleted/config"
	modeldb "github.com/scch94/agentsDeleted/models/db"
	"github.com/scch94/ins_log"
)

const (
	postgresGetUsers = "SELECT ua.oid, ua.user_id, cu.client_oid FROM user_adm ua JOIN client_user cu ON ua.oid = cu.user_oid WHERE cu.client_oid = $1 AND cu.tenant_oid=$2"
)

func GetUsers(agentOid string, ctx context.Context) ([]modeldb.UsersDb, error) {
	//establece el contexto actual
	ctx = ins_log.SetPackageNameInContext(ctx, "database")

	//creamos una lista donde guardaremos el resultado de la operaicon!
	var users []modeldb.UsersDb

	var err error = nil

	ins_log.Tracef(ctx, "starting to get the Users por the agent with oid :%v", agentOid)

	startTime := time.Now()

	ins_log.Tracef(ctx, "this is the QUERY: %s and the params: agentID=%s, and tenant_oid=%s", postgresGetUsers, agentOid, config.Config.Tenant)

	db := GetDb()

	rows, err := db.QueryContext(ctx, postgresGetUsers, agentOid, config.Config.Tenant)
	if err != nil {
		ins_log.Errorf(ctx, "query error %v", err)
		return nil, err
	}
	defer rows.Close()

	// Iterar sobre las filas de resultados
	for rows.Next() {
		var user modeldb.UsersDb
		//inicialisamos el candelete en true
		user.CanDelete = true
		err = rows.Scan(&user.UserOid, &user.UserId, &user.AgentOid)
		if err != nil {
			ins_log.Errorf(ctx, "error scanning row: %v", err)
			return nil, err
		}
		users = append(users, user)
	}
	// Verificar si hubo errores en el procesamiento de las filas
	if err = rows.Err(); err != nil {
		ins_log.Errorf(ctx, "error during row iteration: %v", err)
		return nil, err
	}
	duration := time.Since(startTime)
	ins_log.Infof(ctx, "the query in the database tooks: %v", duration)

	return users, nil

}
