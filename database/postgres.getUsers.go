package database

import (
	"context"
	"fmt"
	"time"

	"github.com/scch94/agentsDeleted/config"
	modeldb "github.com/scch94/agentsDeleted/models/db"
	modelUtils "github.com/scch94/agentsDeleted/models/utils"
	"github.com/scch94/ins_log"
)

const (
	postgresGetUsers = "SELECT ua.oid, ua.user_id, cu.client_oid FROM user_adm ua FULL JOIN client_user cu ON ua.oid = cu.user_oid WHERE cu.client_oid = $1 AND cu.tenant_oid=$2"
)

func GetUsers(ctx context.Context, agent *modelUtils.Agents) ([]modeldb.UsersDb, error) {
	//establece el contexto actual
	ctx = ins_log.SetPackageNameInContext(ctx, "database")

	//creamos una lista donde guardaremos el resultado de la operaicon! ademas inciale
	var users []modeldb.UsersDb

	var err error = nil

	ins_log.Tracef(ctx, "starting to get the Users por the agent with oid :%v", agent.AgentOid)

	startTime := time.Now()

	ins_log.Tracef(ctx, "this is the QUERY: %s and the params: agentID=%s, and tenant_oid=%s", postgresGetUsers, agent.AgentOid, config.Config.Tenant)

	db := GetDb()

	rows, err := db.QueryContext(ctx, postgresGetUsers, agent.AgentOid, config.Config.Tenant)
	if err != nil {
		ins_log.Errorf(ctx, "query error %v", err)
		return nil, err
	}
	defer rows.Close()

	//tiempo de respuesta
	duration := time.Since(startTime)
	ins_log.Infof(ctx, "the query in the database tooks: %v", duration)

	// Iterar sobre las filas de resultados
	for rows.Next() {
		var user modeldb.UsersDb
		err = rows.Scan(&user.UserOid, &user.UserId, &user.AgentOid)
		if err != nil {
			ins_log.Errorf(ctx, "error scanning row: %v", err)
			return nil, err
		}
		//chequeamos si el usuario puede ser elimiando
		err = isUserParent(ctx, &user)
		if err != nil {
			ins_log.Errorf(ctx, "error in the function is UserParent err : %v ", err)
			return nil, err
		}
		if !user.CanDelete.AgentcanDeleted {
			if agent.CanDelete.Reason == "" {
				agent.CanDelete.Reason = fmt.Sprintf("agent whit id %v, can not deleted reasons :", agent.AgentId)
			}
			agent.CanDelete.AgentcanDeleted = user.CanDelete.AgentcanDeleted
			agent.CanDelete.Reason = fmt.Sprintf("%v %v.", agent.CanDelete.Reason, user.CanDelete.Reason)
		}
		users = append(users, user)

	}
	// Verificar si hubo errores en el procesamiento de las filas
	if err = rows.Err(); err != nil {
		ins_log.Errorf(ctx, "error during row iteration: %v", err)
		return nil, err
	}

	return users, nil

}
