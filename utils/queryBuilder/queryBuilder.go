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
	clientOperative                     = "client_operative"

	//agent table to delete
	agentTableName      = "agent"
	accountTableName    = "account"
	agentClerkTableName = "agent_clerk"
	agentZoneTableName  = "agent_zone"
	agentChild          = "agent_child"

	//conditionals
	oid       = "oid"
	clientOid = "client_oid"
	agentOid  = "agent_oid"
	childOid  = "CHILD_OID"
)

// CreateQuery construye una consulta SQL para eliminar datos de una tabla dada con condiciones específicas.
func CreateQuery(ctx context.Context, tableName string, conditional string, infodb []modeldb.ModelsDb) string {
	// Consulta base para eliminación
	const deleteQueryTemplate = "DELETE FROM %s WHERE %s IN ("

	// Establece el contexto para los registros
	ctx = ins_log.SetPackageNameInContext(ctx, "queryBuilder")
	ins_log.Infof(ctx, "Iniciando la creación de una consulta para eliminar datos de la tabla %s", tableName)

	// Forma la consulta base
	deleteQuery := fmt.Sprintf(deleteQueryTemplate, tableName, conditional)

	// Usa strings.Builder para construir la consulta
	var sb strings.Builder
	sb.WriteString(deleteQuery)

	// Si no hay datos para eliminar, devuelve una cadena vacía y registra un error
	if len(infodb) == 0 {
		ins_log.Tracef(ctx, "No hay datos para eliminar")
		return ""
	}

	// Agrega las condiciones a la consulta
	isFirst := true
	for _, item := range infodb {
		if item.CanDeleted() && item.Condition() != "" {
			if !isFirst {
				sb.WriteString(", ")
			}
			sb.WriteString(item.Condition())
			isFirst = false
		}
	}

	// Cierra el paréntesis y añade un salto de línea
	sb.WriteString(");\n")
	return sb.String()
}

func AgentQueyBuilders(ctx context.Context, agents []modelUtils.Agents) []modelUtils.Table {

	//esta consulta despuesdebera ir en constantes
	var deleteQuery = "delete from %v where %v in ("
	flag := false

	ctx = ins_log.SetPackageNameInContext(ctx, "queryBuilder")
	ins_log.Infof(ctx, "starting to create the querys to delete agents")

	//creamos un arreglo con el nombre de las tablas a eliminar
	var tablesToDelete []modelUtils.Table = []modelUtils.Table{
		{TableName: accountTableName, Conditional: agentOid, QueryToDelete: strings.Builder{}},
		{TableName: agentChild, Conditional: childOid, QueryToDelete: strings.Builder{}},
		{TableName: agentClerkTableName, Conditional: agentOid, QueryToDelete: strings.Builder{}},
		{TableName: agentZoneTableName, Conditional: agentOid, QueryToDelete: strings.Builder{}},
		{TableName: agentTableName, Conditional: oid, QueryToDelete: strings.Builder{}},
		{TableName: clientOperative, Conditional: clientOid, QueryToDelete: strings.Builder{}},
		{TableName: clientProfileClientComTableName, Conditional: clientOid, QueryToDelete: strings.Builder{}},
		{TableName: clientProfileClientNotifTableName, Conditional: clientOid, QueryToDelete: strings.Builder{}},
		{TableName: clientProfileClienRestrictTableName, Conditional: clientOid, QueryToDelete: strings.Builder{}},
		{TableName: clientServiceInterfaceTableName, Conditional: clientOid, QueryToDelete: strings.Builder{}},
		{TableName: clientTableName, Conditional: oid, QueryToDelete: strings.Builder{}},
	}

	for i := 0; i < len(agents); i++ {
		if agents[i].CanDelete.AgentcanDeleted {

			flag = true
			//sobre cada agente debembemos rrecorres el tables to delete para crear el query de cada tabla
			for j := 0; j < len(tablesToDelete); j++ {
				//si es la primera ves arrancamos la query si no le agregamos una coma
				if i == 0 || tablesToDelete[j].QueryToDelete.Len() == 0 {
					startOfTheQuery := fmt.Sprintf(deleteQuery, tablesToDelete[j].TableName, tablesToDelete[j].Conditional)
					tablesToDelete[j].QueryToDelete.WriteString(startOfTheQuery)

				} else {
					tablesToDelete[j].QueryToDelete.WriteString(",")
				}
				tablesToDelete[j].QueryToDelete.WriteString(agents[i].AgentOid)
			}
		}

	}

	// Cerramos las consultas
	for i := range tablesToDelete {
		if !flag {
			tablesToDelete[i].QueryToDelete.WriteString("no agents to delete ();\n")
		} else {
			tablesToDelete[i].QueryToDelete.WriteString(");\n")
		}

	}
	return tablesToDelete
}
