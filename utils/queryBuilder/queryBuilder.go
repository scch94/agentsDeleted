package querybuilder

import (
	"context"
	"fmt"
	"strings"

	modeldb "github.com/scch94/agentsDeleted/models/db"
	modelUtils "github.com/scch94/agentsDeleted/models/utils"
	"github.com/scch94/ins_log"
)

const (
	//client tables to delete
	clientTableName                     = "client"
	clientProfileClientComTableName     = "client_profile_client_comm"
	clientProfileClientNotifTableName   = "client_profile_client_notif"
	clientProfileClienRestrictTableName = "client_profile_client_restric"
	clientServiceInterfaceTableName     = "client_service_interface"

	//agent table to delete
	agentTableName      = "agent"
	accountTableName    = "account"
	agentClerkTableName = "agent_clerk"
	agentZoneTableName  = "agent_zone"

	//conditionals
	oid       = "oid"
	clientOid = "client_oid"
	agentOid  = "agent_oid"
)

func CreateQuery(ctx context.Context, tableName string, conditional string, infodb []modeldb.ModelsDb) string {
	//esta consulta despuesdebera ir en constantes
	var deleteQuery = "delete from %v where %v in ("

	ctx = ins_log.SetPackageNameInContext(ctx, "queryBuilder")
	ins_log.Infof(ctx, "starting to create a query to delete data to the table %v", tableName)

	deleteQuery = fmt.Sprintf(deleteQuery, tableName, conditional)

	//creamos strings builder para manejar la consulta
	var sb strings.Builder
	sb.WriteString(deleteQuery)

	//si no tiene data en el arreglo se devuelve un error
	if len(infodb) == 0 {
		ins_log.Errorf(ctx, "no data to delete")
		return ""
	}
	for i, info := range infodb {

		if info.CanDeleted() {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(info.Condition())
		}

	}

	//agregamos el parentesis al final
	sb.WriteString(");\n")
	return sb.String()
}

func AgentQueyBuilders(ctx context.Context, agents []modelUtils.Agents) []modelUtils.Table {

	//esta consulta despuesdebera ir en constantes
	var deleteQuery = "delete from %v where %v in ("

	ctx = ins_log.SetPackageNameInContext(ctx, "queryBuilder")
	ins_log.Infof(ctx, "starting to create the querys to delete agents")

	//creamos un arreglo con el nombre de las tablas a eliminar
	var tablesToDelete []modelUtils.Table = []modelUtils.Table{
		{TableName: clientTableName, Conditional: oid, QueryToDelete: strings.Builder{}},
		{TableName: clientProfileClientComTableName, Conditional: clientOid, QueryToDelete: strings.Builder{}},
		{TableName: clientProfileClientNotifTableName, Conditional: clientOid, QueryToDelete: strings.Builder{}},
		{TableName: clientProfileClienRestrictTableName, Conditional: clientOid, QueryToDelete: strings.Builder{}},
		{TableName: clientServiceInterfaceTableName, Conditional: clientOid, QueryToDelete: strings.Builder{}},
		{TableName: agentTableName, Conditional: oid, QueryToDelete: strings.Builder{}},
		{TableName: accountTableName, Conditional: agentOid, QueryToDelete: strings.Builder{}},
		{TableName: agentClerkTableName, Conditional: agentOid, QueryToDelete: strings.Builder{}},
		{TableName: agentZoneTableName, Conditional: agentOid, QueryToDelete: strings.Builder{}},
	}
	for i := 0; i < len(agents); i++ {
		if agents[i].CanDelete.AgentcanDeleted {
			//sobre cada agente debembemos rrecorres el tables to delete para crear el query de cada tabla
			for j := 0; j < len(tablesToDelete); j++ {
				//si es la primera ves arrancamos la query si no le agregamos una coma
				if i == 0 {
					startOfTheQuery := fmt.Sprintf(deleteQuery, tablesToDelete[j].TableName, tablesToDelete[j].Conditional)
					tablesToDelete[j].QueryToDelete.WriteString(startOfTheQuery)
				} else {
					tablesToDelete[j].QueryToDelete.WriteString(",")
				}

				tablesToDelete[j].QueryToDelete.WriteString(agents[i].AgentOid)
				//si es el ultimo agente cierro el corchete y pongo el ;
				if i == len(agents)-1 {
					tablesToDelete[j].QueryToDelete.WriteString(");\n")
				}
			}
		}
	}
	return tablesToDelete
}
