package modelUtils

type Agents struct {
	AgentId   string
	AgentOid  string
	Credit    string
	CanDelete canDeleted
}

type canDeleted struct {
	AgentcanDeleted bool
	Reason          string
}
