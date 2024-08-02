package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/scch94/agentsDeleted/config"
	modelUtils "github.com/scch94/agentsDeleted/models/utils"
	"github.com/scch94/ins_log"
)

const (
	oracleIsagentParent = "SELECT * FROM agent WHERE parent_oid = :1 AND tenant_oid = :2 AND ROWNUM = 1"
)

func IsAgentParentOracle(ctx context.Context, agent *modelUtils.Agents) error {

	// Establece el contexto actual
	ctx = ins_log.SetPackageNameInContext(ctx, "database")
	ins_log.Tracef(ctx, "starting to check if the Agent with id %v is a parent", agent.AgentId)

	// Creamos strings builder para manejar el texto
	var sb strings.Builder
	sb.WriteString(agent.CanDelete.Reason)

	// Si el agente no tiene un texto para ser eliminado, inicializamos el texto para el log final
	if agent.CanDelete.Reason == "" {
		sb.WriteString(fmt.Sprintf("agent with id %v, cannot be deleted for reasons:", agent.AgentId))
	}

	//inicamos el contador de la consulta a la base !
	startTime := time.Now()

	ins_log.Tracef(ctx, "this is the QUERY: %s and the params: agentID=%s, and tenant_oid=%v", oracleIsagentParent, agent.AgentOid, config.Config.Tenant)

	db := GetDb()

	rows, err := db.QueryContext(ctx, oracleIsagentParent, agent.AgentOid, config.Config.Tenant)
	if err != nil {
		ins_log.Errorf(ctx, "query error %v", err)
		return err
	}

	if rows.Next() {
		ins_log.Infof(ctx, "Agent %v is a parent and cannot be deleted", agent.AgentOid)
		sb.WriteString(" Agent is a parent and cannot be deleted.")
		agent.CanDelete.AgentcanDeleted = false
		agent.CanDelete.Reason = sb.String()

	} else {
		ins_log.Infof(ctx, "Agent %v is not a parent and can be deleted", agent.AgentOid)
	}
	// Verificar si hubo errores en el procesamiento de las filas
	if err = rows.Err(); err != nil {
		ins_log.Errorf(ctx, "error during row iteration: %v", err)
		return err
	}
	duration := time.Since(startTime)
	ins_log.Infof(ctx, "the query in the database tooks: %v", duration)

	return nil

}
