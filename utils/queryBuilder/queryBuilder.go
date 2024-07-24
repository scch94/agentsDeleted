package querybuilder

import (
	"context"
	"fmt"
	"strings"

	modeldb "github.com/scch94/agentsDeleted/models/db"
	"github.com/scch94/ins_log"
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
