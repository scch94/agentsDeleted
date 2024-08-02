package config

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/scch94/Gconfiguration"
	"github.com/scch94/ins_log"
)

////EPIN_NEW/epin@192.168.0.157:1521/epin

var Config AgentsToDeletedConfiguration

type AgentsToDeletedConfiguration struct {
	LogLevel                 string `json:"log_level"`
	Log_name                 string `json:"log_name"`
	DatabaseConnectionString string `json:"database_connection_string"`
	Tenant                   int    `json:"tenant"`
	UbSicationAgentFile      string `json:"ubication_agents_file"`
	DatabaseEngine           string `json:"database_engine"`
}

func Upconfig(ctx context.Context) error {
	//traemos el contexto y le setiamos el contexto actual
	ctx = ins_log.SetPackageNameInContext(ctx, "config")

	ins_log.Info(ctx, "starting to get the config struct ")
	err := Gconfiguration.GetConfig(&Config, "../config", "agentsDeleted.json")
	if err != nil {
		ins_log.Fatalf(ctx, "error in Gconfiguration.GetConfig() ", err)
		return err
	}
	return nil
}
func (a AgentsToDeletedConfiguration) ConfigurationString() string {
	configJSON, err := json.Marshal(a)
	if err != nil {
		return fmt.Sprintf("Error al convertir la configuración a JSON: %v", err)
	}
	return string(configJSON)
}
