package controller

import (
	"context"
	"fmt"
	"os"

	"github.com/scch94/agentsDeleted/database"
	modeldb "github.com/scch94/agentsDeleted/models/db"
	modelUtils "github.com/scch94/agentsDeleted/models/utils"
	fileWriter "github.com/scch94/agentsDeleted/utils/filewriter"
	querybuilder "github.com/scch94/agentsDeleted/utils/queryBuilder"
	"github.com/scch94/ins_log"
)

func DeleteMsisdnAgents(ctx context.Context, agents *[]modelUtils.Agents) error {
	ctx = ins_log.SetPackageNameInContext(ctx, "controller")
	ins_log.Infof(ctx, "starting to create the script to delete agents_mobile")

	//llamamos la funcion que nos traera la informacion de cada movil vincualdo a los agentes que debemos elimar
	msisdnsInfo, err := getAgentsInfo(ctx, agents)
	if err != nil {
		ins_log.Errorf(ctx, "error on the function getAgentsInfo() err: %v", err)
		return err
	}
	var modelsDbInfo []modeldb.ModelsDb
	for _, msisdn := range msisdnsInfo {
		modelsDbInfo = append(modelsDbInfo, &msisdn)
	}

	//primero eliminaremos el agent_movil_pin
	queryToDeleteAgentMobilePin := querybuilder.CreateQuery2(ctx, "agent_mobile_pin", "msisdn_oid", modelsDbInfo)
	err = fileWriter.WriteInAfile(ctx, queryToDeleteAgentMobilePin, "../scripts/agent_mobile_scripts.txt", "query to delete msisdn_pin for an agents in the list")
	if err != nil {
		ins_log.Errorf(ctx, "error when we try to write in a file the queryToDeleteAgentMobilePin and the error message is: %s", err)
		return err
	}
	ins_log.Infof(ctx, "query to delete agent mobile pin was created and writed")

	//ahora con esa lista crearemos el archivo con el script para eliminar el agent movil de los agentes
	queryToDeleteAgentMobile := querybuilder.CreateQuery2(ctx, "agent_mobile", "oid", modelsDbInfo)
	//pasamos el texto de la query y la ubicacion del archivo
	err = fileWriter.WriteInAfile(ctx, queryToDeleteAgentMobile, "../scripts/agent_mobile_scripts.txt", "query to delete msisdn for an agents in the list")
	if err != nil {
		ins_log.Errorf(ctx, "error when we try to write in a file the querytodeletemsisdn and the error message is: %s", err)
		return err
	}
	ins_log.Infof(ctx, "query to delete agent mobile was created and writed")
	return nil
}

func getAgentsInfo(ctx context.Context, agents *[]modelUtils.Agents) ([]modeldb.MsisdnDb, error) {

	ins_log.Tracef(ctx, "starting to get agents_msisdn info")

	// Creamos una lista para almacenar los números MSISDN
	var AllmsisdnsInfo []modeldb.MsisdnDb

	// Abre el archivo en modo de escritura (crea el archivo si no existe)
	file, err := os.OpenFile("../utils/agents_mobile.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		ins_log.Errorf(ctx, "error opening file: %v", err)
		return nil, err
	}
	defer file.Close()

	//abrimos el archivo donde loguearemos los moviles que se eliminaran
	for i := range *agents {
		msisdnsInfo, err := database.GetMsisdn(ctx, &(*agents)[i])
		if err != nil {
			ins_log.Errorf(ctx, "error in getmsisdn: %v", err)
			return nil, err
		}
		ins_log.Tracef(ctx, "the agent with id %v has %v moviles", (*agents)[i].AgentId, len(msisdnsInfo))
		for _, msisdnInfo := range msisdnsInfo {
			if _, err := file.WriteString(fmt.Sprintf("Agent ID: %v, Mobile to delete: %v\n", (*agents)[i].AgentId, msisdnInfo.Msisdn)); err != nil {
				ins_log.Errorf(ctx, "error writing to file: %v", err)
				return nil, err
			}
		}
		AllmsisdnsInfo = append(AllmsisdnsInfo, msisdnsInfo...)

	}
	return AllmsisdnsInfo, nil
}
