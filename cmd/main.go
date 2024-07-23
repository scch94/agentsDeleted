package main

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/scch94/agentsDeleted/config"
	"github.com/scch94/agentsDeleted/database"
	controller "github.com/scch94/agentsDeleted/utils/controllers"
	reader "github.com/scch94/agentsDeleted/utils/filereader"
	"github.com/scch94/ins_log"
)

func main() {

	//creamos contecxto par ala ejecucion
	ctx := context.Background()

	//creamos el archivo de donde el programa logueara
	go initializeAndWatchLogger(ctx)

	//levantamos la configuracion
	if err := config.Upconfig(ctx); err != nil {
		ins_log.Errorf(ctx, "error loading configuration: %v", err)
		return
	}

	//inicamos el logger
	ins_log.SetService("agentsdeleter")
	ins_log.SetLevel(config.Config.LogLevel)
	ctx = ins_log.SetPackageNameInContext(ctx, "main")

	//incialisamos la base de datos
	database.InitDb(ctx)

	//leemoos y recuperamos la lista de agentes
	agents, err := reader.Read(ctx)
	if err != nil {
		ins_log.Fatalf(ctx, "error triyig to read the agents file: %v", err)

	}
	ins_log.Infof(ctx, "this are the agents in the file %v", agents)

	//llamamos la funcion para crear el script para eliminar los moviles de los agentes en la base
	err = controller.DeleteMsisdnAgents(ctx, agents)
	if err != nil {
		ins_log.Errorf(ctx, "error creating the script to delete the msisdns of the agents: %v", err)
		return
	}

	//llamamos la funcion para crear el script para eliminar los users de los agentes en la base
	err = controller.DeleteUserAdm(ctx, agents)
	if err != nil {
		ins_log.Errorf(ctx, "error creating the script to delete the users of the agents: %v", err)
		return
	}

	// ins_log.Infof(ctx, "starting to check the users vinculated to the agents")

	// users, err := getUser_adm(msisdns, ctx)
	// if err != nil {
	// 	ins_log.Errorf(ctx, "error when we try to getUser:adm(): %v", err)
	// 	return
	// }
	// ins_log.Debugf(ctx, "Users: %v", users)
}
func initializeAndWatchLogger(ctx context.Context) {
	var file *os.File
	var logFileName string
	var err error
	for {
		select {
		case <-ctx.Done():
			return
		default:
			logDir := "../log"

			// Create the log directory if it doesn't exist
			if err = os.MkdirAll(logDir, 0755); err != nil {
				ins_log.Errorf(ctx, "error creating log directory: %v", err)
				return
			}

			// Define the log file name
			today := time.Now().Format("2006-01-02 15")
			replacer := strings.NewReplacer(" ", "_")
			today = replacer.Replace(today)
			logFileName = filepath.Join(logDir, config.Config.Log_name+today+".log")

			// Open the log file
			file, err = os.OpenFile(logFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				ins_log.Errorf(ctx, "error opening log file: %v", err)
				return
			}

			// Create a writer that writes to both file and console
			multiWriter := io.MultiWriter(os.Stdout, file)
			ins_log.StartLoggerWithWriter(multiWriter)

			// Esperar hasta el inicio de la próxima hora
			nextHour := time.Now().Truncate(time.Hour).Add(time.Hour)
			time.Sleep(time.Until(nextHour))

			// Close the previous log file
			file.Close()
		}
	}
}
