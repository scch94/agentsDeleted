package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/scch94/agentsDeleted/config"
	modeldb "github.com/scch94/agentsDeleted/models/db"
	modelUtils "github.com/scch94/agentsDeleted/models/utils"
	"github.com/scch94/ins_log"
)

const (
	oracleGetMsisdn = `SELECT AM.OID, A.OID, MSISDN, ACC.CREDIT FROM agent A FULL JOIN agent_mobile AM ON AM.AGENT_OID=A.OID FULL JOIN account ACC ON ACC.AGENT_OID = A.OID WHERE A.AGENT_ID=:1 AND A.TENANT_OID=:2`
)

func GetMsisdnOracle(ctx context.Context, agent *modelUtils.Agents) ([]modeldb.MsisdnDb, error) {
	//estable el contexto acutual
	// Establece el contexto actual
	ctx = ins_log.SetPackageNameInContext(ctx, "database")

	// Creamos una lista para almacenar los números MSISDN
	var msisdnsInfo []modeldb.MsisdnDb

	ins_log.Tracef(ctx, "starting to get the msisdn por the agent with id :%v", agent.AgentId)
	startTime := time.Now()

	ins_log.Tracef(ctx, "this is the QUERY: %s and the params: agentID=%s, and tenant_oid=%v", oracleGetMsisdn, agent.AgentId, config.Config.Tenant)

	db := GetDb()

	rows, err := db.QueryContext(ctx, oracleGetMsisdn, agent.AgentId, config.Config.Tenant)
	if err != nil {
		ins_log.Errorf(ctx, "query error %v", err)
		return nil, err
	}
	defer rows.Close()
	// Iterar sobre las filas de resultados
	for rows.Next() {
		var msisdnSql modeldb.MsisdnDbSql
		var credit sql.NullFloat64
		err := rows.Scan(&msisdnSql.MsisdnOid, &msisdnSql.AgentOid, &msisdnSql.Msisdn, &credit)
		if err != nil {
			ins_log.Errorf(ctx, "error scanning row: %v", err)
			return nil, err
		}
		ins_log.Infof(ctx, "estos son los resultados de la consulta ! 1:%v   2:%v 	3:%v	4:%v", msisdnSql.MsisdnOid, msisdnSql.AgentOid, msisdnSql.Msisdn, credit)
		msisdnInfo := msisdnSql.ConvertMsisdn()
		msisdnsInfo = append(msisdnsInfo, msisdnInfo)
		agent.AgentOid = msisdnInfo.AgentOid
		agent.Credit = credit.Float64
		ins_log.Infof(ctx, "esto solo es para ver la info recuperada ")
	}

	// Verificar si hubo errores en el procesamiento de las filas
	if err = rows.Err(); err != nil {
		ins_log.Errorf(ctx, "error during row iteration: %v", err)
		return nil, err
	}
	duration := time.Since(startTime)
	ins_log.Infof(ctx, "the query in the database tooks: %v", duration)

	ins_log.Infof(ctx, "msisdnsinfos %v", msisdnsInfo)
	ins_log.Infof(ctx, "AGENTS: %v", agent.AgentOid)

	return msisdnsInfo, nil

}
