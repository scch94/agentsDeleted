package modelUtils

type Agents struct {
	AgentId   string
	AgentOid  string
	Credit    float64
	CanDelete canDeleted
}

type canDeleted struct {
	AgentcanDeleted bool
	Reason          string
}
