package modelUtils

type Agents struct {
	AgentId   string
	AgentOid  string
	CanDelete canDeleted
}

type canDeleted struct {
	AgentcanDeleted bool
	Reason          string
}
