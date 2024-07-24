package reader

import (
	"bufio"
	"context"
	"os"

	"github.com/scch94/agentsDeleted/config"
	modelUtils "github.com/scch94/agentsDeleted/models/utils"
	"github.com/scch94/ins_log"
)

func Read(ctx context.Context) ([]modelUtils.Agents, error) {
	ctx = ins_log.SetPackageNameInContext(ctx, "fileReader")

	//traemos del config la ubicaicon del archivo de eliminar agentes. y creamos la variable donde guardaremos el agent_id del archivo
	ubication := config.Config.UbSicationAgentFile
	var agents []modelUtils.Agents

	ins_log.Infof(ctx, "starting to get the agents  ubication of the file is %v", ubication)

	//abrimos el archivo lo leemos y hacemos un defer para que lo cierre al final de la funcion ,
	file, err := os.Open(ubication)
	if err != nil {
		ins_log.Fatalf(ctx, "error opening the agent file: %v", err)
		return nil, err
	}
	defer file.Close()

	//scaneamos
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var agent modelUtils.Agents
		agent.AgentId = scanner.Text()
		agent.CanDelete.AgentcanDeleted = true
		agent.CanDelete.Reason = ""

		agents = append(agents, agent)
	}

	//verificamos si hubo algun error en el scaneo
	if err := scanner.Err(); err != nil {
		ins_log.Fatalf(ctx, "error scanning the file: %v", err)
		return nil, err
	}

	return agents, nil

}
