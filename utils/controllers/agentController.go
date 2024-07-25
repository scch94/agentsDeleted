package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/scch94/agentsDeleted/database"
	modelUtils "github.com/scch94/agentsDeleted/models/utils"
	fileWriter "github.com/scch94/agentsDeleted/utils/filewriter"
	"github.com/scch94/ins_log"
)

func DeleteAgents(ctx context.Context, agents *[]modelUtils.Agents) error {
	ctx = ins_log.SetPackageNameInContext(ctx, "controller")
	ins_log.Infof(ctx, "starting to create the script to delete agents ")

	//chequeamos si el saldo del usuario es =0
	checkCreditAccount(ctx, agents)
	ins_log.Infof(ctx, "check credit accounts finish")

	//veremos si el agente es padre de algun otro agente o si tiene vista de algun otro agente
	err := isAgentsParent(ctx, agents)
	if err != nil {
		ins_log.Errorf(ctx, "error in the function isAgentsParent() err: %v", err)
		return err
	}
	ins_log.Infof(ctx, "check if the agents are parents finish")

	//falta la creacion de scripts para elimianr agentes.

	//creamos el archivo con concluisones
	conclusionText := ConclusionTextBuilder(ctx, agents)
	err = fileWriter.WriteInAfile(ctx, conclusionText, "../utils/conclusions.txt", "Conclusions of the scriptBuilder")
	if err != nil {
		ins_log.Errorf(ctx, "error when we try to write in a file the conclusionText and the error message is: %s", err)
		return err
	}
	ins_log.Infof(ctx, "conclusion file was created ")

	return nil
}

func ConclusionTextBuilder(ctx context.Context, agents *[]modelUtils.Agents) string {
	//creamos strings builder para manejar el texto
	var sb strings.Builder

	for i, agent := range *agents {
		sb.WriteString("\n")
		if agent.CanDelete.AgentcanDeleted {
			sb.WriteString(fmt.Sprintf("%v the agent whit id %v can be delete", i, agent.AgentId))
		} else {
			sb.WriteString(fmt.Sprintf("%v - %v", i, agent.CanDelete.Reason))
		}
	}
	return sb.String()
}

// funcion para mirar que agente de la lista no puede ser eliminado por que tiene saldo mayor a 0
func checkCreditAccount(ctx context.Context, agents *[]modelUtils.Agents) {

	ins_log.Tracef(ctx, "checking if the agents have a credit different than 0")

	for i := range *agents {
		//creamos strings builder para manejar el texto
		var sb strings.Builder
		sb.WriteString((*agents)[i].CanDelete.Reason)
		if (*agents)[i].Credit != "0" {
			(*agents)[i].CanDelete.AgentcanDeleted = false
			sb.WriteString(" Agent has a credit different than 0. ")
			(*agents)[i].CanDelete.Reason = sb.String()
		}

	}

}

func isAgentsParent(ctx context.Context, agents *[]modelUtils.Agents) error {

	ins_log.Tracef(ctx, "checking if the agents are parents")

	for i := range *agents {
		err := database.IsAgentParent(ctx, &(*agents)[i])
		if err != nil {
			ins_log.Errorf(ctx, "error in the database function IsAgentParent err: %v", err)
			return err
		}
	}
	return nil
}
