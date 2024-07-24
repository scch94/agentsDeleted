package controller

import (
	"context"

	modelUtils "github.com/scch94/agentsDeleted/models/utils"
	"github.com/scch94/ins_log"
)

func DeleteAgents(ctx context.Context, agents *[]modelUtils.Agents) error {
	ctx = ins_log.SetPackageNameInContext(ctx, "controller")
	ins_log.Infof(ctx, "starting to create the script to delete agents : %v", agents)

	return nil
}
