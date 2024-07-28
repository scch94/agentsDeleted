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

func DeleteUserAdm(ctx context.Context, agents *[]modelUtils.Agents) error {
	ctx = ins_log.SetPackageNameInContext(ctx, "controller")
	ins_log.Infof(ctx, "starting to create the script to delete Users vinculated to the agents in the list")

	//llamamos la funcion que nos traera la informacion de cada movil vinculado a los agentes que debemos eliminar
	usersInfo, err := getUsersAdm(ctx, agents)
	if err != nil {
		ins_log.Errorf(ctx, "error on the function getUsersAdm() err: %v", err)
		return err
	}
	ins_log.Infof(ctx, "usuarios : %v agentes %v", usersInfo, agents)

	var modelsDbInfo []modeldb.ModelsDb
	for _, user := range usersInfo {
		modelsDbInfo = append(modelsDbInfo, &user)
	}
	if len(usersInfo) > 0 {
		//eliminamos la old_password del user
		queryToDeleteUserOldPassword := querybuilder.CreateQuery(ctx, "old_password", "user_oid", modelsDbInfo)
		//pasamos el texto de la query y la ubicacion del archivo para crear el script
		err = fileWriter.WriteInAfile(ctx, queryToDeleteUserOldPassword, "../scripts/users.txt", "query to delete the old password for the users vinculated to the agents in the list")
		if err != nil {
			ins_log.Errorf(ctx, "error when we try to write in a file the queryToDeleteUserOldPassword and the error message is: %s", err)
			return err
		}
		ins_log.Infof(ctx, "query to delete users_old_password was created and writed")

		//ahora vamos a eliminar el client_user
		querytoDeleteClientUser := querybuilder.CreateQuery(ctx, "client_user", "user_oid", modelsDbInfo)
		//pasamos el texto de la query y la ubicacion del archivo para crear el script
		err = fileWriter.WriteInAfile(ctx, querytoDeleteClientUser, "../scripts/users.txt", "query to delete the client_user for the users vinculated to the agents in the list")
		if err != nil {
			ins_log.Errorf(ctx, "error when we try to write in a file the querytoDeleteClientUser and the error message is: %s", err)
			return err
		}
		ins_log.Infof(ctx, "query to delete ClientUser was created and writed")

		//ahora vamos a eliminar el user_adm
		queryToDeleteUserAdm := querybuilder.CreateQuery(ctx, "user_adm", "oid", modelsDbInfo)
		//pasamos el texto de la query y la ubicacion del archivo para crear el script
		err = fileWriter.WriteInAfile(ctx, queryToDeleteUserAdm, "../scripts/users.txt", "query to delete the user_adm for the users vinculated to the agents in the list")
		if err != nil {
			ins_log.Errorf(ctx, "error when we try to write in a file the queryToDeleteUserAdm and the error message is: %s", err)
			return err
		}
		ins_log.Infof(ctx, "query to delete user_adm was created and writed")
	} else {
		ins_log.Infof(ctx, "No users to delete")
	}

	return nil

}

func getUsersAdm(ctx context.Context, agents *[]modelUtils.Agents) ([]modeldb.UsersDb, error) {

	ins_log.Tracef(ctx, "starting to get users")

	//creamos una lista para almacenar los usuarios
	var AllUsers []modeldb.UsersDb

	// Abre el archivo en modo de escritura (crea el archivo si no existe)
	file, err := os.OpenFile("../utils/userAdm.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		ins_log.Errorf(ctx, "error opening file: %v", err)
		return nil, err
	}
	defer file.Close()

	//abrimos el archivo donde se loguearan los usarios que se trataran de eliminar.
	ins_log.Infof(ctx, "AGENTS:", agents)
	for i := range *agents {
		usersInfo, err := database.GetUsers(ctx, &(*agents)[i])
		if err != nil {
			ins_log.Errorf(ctx, "error in getUsers: %v", err)

		}
		ins_log.Tracef(ctx, "the agent with id %v has %v users", (*agents)[i].AgentId, len(usersInfo))
		for _, userInfo := range usersInfo {
			if _, err := file.WriteString(fmt.Sprintf("agent ID: %v, user: %v can delete %v\n", (*agents)[i].AgentId, userInfo.UserId, userInfo.CanDelete.AgentcanDeleted)); err != nil {
				ins_log.Errorf(ctx, "error writing to file: %v", err)
			}
		}
		AllUsers = append(AllUsers, usersInfo...)
	}
	return AllUsers, nil
}
