package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/scch94/agentsDeleted/config"
	"github.com/scch94/agentsDeleted/database"
	modeldb "github.com/scch94/agentsDeleted/models/db"
	modelUtils "github.com/scch94/agentsDeleted/models/utils"
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

	//iniciamos el proceso
	agents, err := startProccess(ctx)

	if err != nil {
		ins_log.Fatalf(ctx, "error triyig to read the agents file: %v", err)

	}
	ins_log.Infof(ctx, "this are the agents in the file %v", agents)

	//revisamos si el agente tiene numeros de movil si los tiene crea un archivo txt con el numero de agente y los moviles a
	msisdns, err := getAgentsMobile(agents, ctx)
	if err != nil {
		ins_log.Errorf(ctx, "error in the function getAgentsMobile: %v", err)
		return
	}

	queries, err := BuildSQLQuery(msisdns, ctx)
	if err != nil {
		ins_log.Infof(ctx, "no msisdns to delete")
		return
	}

	ins_log.Infof(ctx, "querytodeletemsisdn is %s", queries.Agent_mobile)
	ins_log.Infof(ctx, "query to delete the pin of the msisdns is %s", queries.Agent_mobile_pin)

	err = writeInAfile(ctx, queries.Agent_mobile, "../utils/agent_mobile_scripts.txt", "query to delete msisdn for an agents in the list")
	if err != nil {
		ins_log.Errorf(ctx, "error when we try to write in a file the querytodeletemsisdn and the error message is: %s", err)
		return
	}
	err = writeInAfile(ctx, queries.Agent_mobile_pin, "../utils/agent_mobile_scripts.txt", "query to delete the pin of the msisdn for an agents in the list")
	if err != nil {
		ins_log.Errorf(ctx, "error when we try to write in a file the querytodeletemsisdn pin and the error message is: %s", err)
		return
	}

	ins_log.Infof(ctx, "starting to check the users vinculated to the agents")

	users, err := getUser_adm(msisdns, ctx)
	if err != nil {
		ins_log.Errorf(ctx, "error when we try to getUser:adm(): %v", err)
		return
	}
	ins_log.Debugf(ctx, "Users: %v", users)

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

const (
	ubicationFile = "../utils/agents.txt"
)

func startProccess(ctx context.Context) ([]int, error) {
	ins_log.Infof(ctx, "starting to get the agents  ubication of the file is %v", ubicationFile)

	file, err := os.Open(ubicationFile)
	if err != nil {
		ins_log.Fatalf(ctx, "error opening the file %v", err)
		return nil, err
	}
	defer file.Close()

	var agents []int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Convertir cada línea a entero
		num, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatalf("Error al convertir el número: %v", err)
			return nil, err
		}
		// Agregar el número al slice
		agents = append(agents, num)
	}
	// Verificar si hubo algún error al leer el archivo
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error al leer el archivo: %v", err)
		return nil, err
	}

	return agents, nil
}

func getUser_adm(agents []modeldb.MsisdnDb, ctx context.Context) ([]modeldb.UsersDb, error) {
	ins_log.Tracef(ctx, "starting to get the users vinculatex to the agents")

	//creamos un slice con todos los usuarios y sus datos que serane limnados
	var allUsers []modeldb.UsersDb

	//creamos el arivo donde escribiremos los usuarios que estan vinculado al agent y lo eliminaremos
	file, err := os.OpenFile("../utils/agent_users.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		ins_log.Errorf(ctx, "error opening file: %v", err)
		return nil, err
	}
	defer file.Close()

	// Mapa para verificar si el oid se repite
	agentsOidMaps := make(map[string]bool)

	//slice con
	var agentsOid []string

	for _, agent := range agents {
		if !agentsOidMaps[agent.AgentOid] {
			agentsOid = append(agentsOid, agent.AgentOid)
			agentsOidMaps[agent.AgentOid] = true
		}
	}

	ins_log.Tracef(ctx, "the agentsOid are %v", agentsOid)

	for _, agentoid := range agentsOid {
		users, err := database.GetUsers(agentoid, ctx)
		if err != nil {
			ins_log.Errorf(ctx, "error getting users(): %v", err)
		}
		ins_log.Infof(ctx, "the agent %v has users: %v", agentoid, len(users))
		for _, user := range users {
			if _, err := file.WriteString(fmt.Sprintf("User to delete: %v", user.UserId)); err != nil {
				ins_log.Errorf(ctx, " error wirting to file %v", err)
				return nil, err
			}
		}
		allUsers = append(allUsers, users...)
	}
	return allUsers, nil
}

func getAgentsMobile(agents []int, ctx context.Context) ([]modeldb.MsisdnDb, error) {
	var AllmsisdnsInfo []modeldb.MsisdnDb
	ins_log.Tracef(ctx, "starting to get agents mobile")
	// Abre el archivo en modo de escritura (crea el archivo si no existe)
	file, err := os.OpenFile("../utils/agents_mobile.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		ins_log.Errorf(ctx, "error opening file: %v", err)
		return nil, err
	}
	defer file.Close()

	for _, agent := range agents {
		msisdnsInfo, err := database.GetMsisdn(agent, ctx)
		if err != nil {
			ins_log.Errorf(ctx, "error getting agent mobiles() ,err: %v", err)
			return nil, err
		}
		ins_log.Infof(ctx, "the agent %v have %v mobiles", agent, len(msisdnsInfo))
		for _, msisdnInfo := range msisdnsInfo {
			if _, err := file.WriteString(fmt.Sprintf("Agent ID: %d, Mobile to delete: %v\n", agent, msisdnInfo.Msisdn)); err != nil {
				ins_log.Errorf(ctx, "error writing to file: %v", err)
				return nil, err
			}
		}
		AllmsisdnsInfo = append(AllmsisdnsInfo, msisdnsInfo...)
	}
	return AllmsisdnsInfo, nil

}

var (
	QUERYTODELETE_AGENT_MOBILE     = "delete from agent_mobile where msisdn in ("
	QUERYTODELETE_AGENT_MOBILE_PIN = "delete from agent_mobile_pin where msisdn_oid in ("
)

func BuildSQLQuery(msisdnsInfo []modeldb.MsisdnDb, ctx context.Context) (modelUtils.FinalQuerys, error) {
	ins_log.Tracef(ctx, "starting to create the query to delete msisdn for and agent")

	//CREAMOS EL STRUCT DONDE GUARDAREMOS LAS QUERYS
	var queries modelUtils.FinalQuerys

	//creamos strings builder para manejar la consulta
	var sb_agent_Mobile strings.Builder
	var sb_agent_Mobile_Pin strings.Builder
	sb_agent_Mobile.WriteString(QUERYTODELETE_AGENT_MOBILE)
	sb_agent_Mobile_Pin.WriteString(QUERYTODELETE_AGENT_MOBILE_PIN)

	//si no tiene numeros a eliminar retorna un string vacio
	if len(msisdnsInfo) == 0 {
		return queries, fmt.Errorf("no msisndns to delete")
	}

	//recorremos el array de msisdn que tenemos
	for i, msisdnInfo := range msisdnsInfo {
		//si no es la primera iteracion le agreagamos una coma para los valores
		if i > 0 {
			sb_agent_Mobile.WriteString(", ")
			sb_agent_Mobile_Pin.WriteString(", ")
		}
		sb_agent_Mobile.WriteString(msisdnInfo.Msisdn)
		sb_agent_Mobile_Pin.WriteString(msisdnInfo.MsisdnOid)
	}

	//agregamos el parentesis al final
	sb_agent_Mobile.WriteString("); \n")
	sb_agent_Mobile_Pin.WriteString("); \n")

	queries.Agent_mobile = sb_agent_Mobile.String()
	queries.Agent_mobile_pin = sb_agent_Mobile_Pin.String()

	return queries, nil
}

func writeInAfile(ctx context.Context, text string, filename string, comment string) error {
	ins_log.Infof(ctx, "starting to write in the file %s", filename)

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		ins_log.Errorf(ctx, "error opening log file: %v", err)
		return err
	}
	defer file.Close()

	comment = "-- " + comment

	_, err = file.WriteString(fmt.Sprintf("%v \n\n%v \n", comment, text))
	if err != nil {
		ins_log.Errorf(ctx, "error writing to file: %v", err)
		return err
	}
	ins_log.Tracef(ctx, "end to wrtite in a file: ")

	return nil
}
