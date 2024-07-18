package database

import (
	"context"
	"time"

	"github.com/scch94/agentsDeleted/config"
	modeldb "github.com/scch94/agentsDeleted/models/db"
	"github.com/scch94/ins_log"
)

const (
	postgresGetMsisdn = `SELECT AM.OID, A.OID, MSISDN FROM AGENT_MOBILE AM INNER JOIN AGENT A ON AM.AGENT_OID=A.OID WHERE A.AGENT_ID=$1 AND A.TENANT_OID=$2`
)

func GetMsisdn(agentID int, ctx context.Context) ([]modeldb.MsisdnDb, error) {
	// Establece el contexto actual
	ctx = ins_log.SetPackageNameInContext(ctx, "database")

	// Creamos una lista para almacenar los números MSISDN
	var msisdnsInfo []modeldb.MsisdnDb

	var err error = nil

	ins_log.Tracef(ctx, "starting to get the msisdn por the agent with id :%v", agentID)
	startTime := time.Now()

	ins_log.Tracef(ctx, "this is the QUERY: %s and the params: agentID=%s, and tenant_oid=%s", postgresGetMsisdn, agentID, config.Config.Tenant)

	db := GetDb()

	rows, err := db.QueryContext(ctx, postgresGetMsisdn, agentID, config.Config.Tenant)
	if err != nil {
		ins_log.Errorf(ctx, "query error %v", err)
		return nil, err
	}
	defer rows.Close()

	// Iterar sobre las filas de resultados
	for rows.Next() {
		var msisdnInfo modeldb.MsisdnDb
		err := rows.Scan(&msisdnInfo.MsisdnOid, &msisdnInfo.AgentOid, &msisdnInfo.Msisdn)
		if err != nil {
			ins_log.Errorf(ctx, "error scanning row: %v", err)
			return nil, err
		}
		msisdnsInfo = append(msisdnsInfo, msisdnInfo)
	}

	// Verificar si hubo errores en el procesamiento de las filas
	if err = rows.Err(); err != nil {
		ins_log.Errorf(ctx, "error during row iteration: %v", err)
		return nil, err
	}
	duration := time.Since(startTime)
	ins_log.Infof(ctx, "the query in the database tooks: %v", duration)

	return msisdnsInfo, nil

}
