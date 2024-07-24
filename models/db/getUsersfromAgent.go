package modeldb

type UsersDb struct {
	UserOid   string
	UserId    string
	AgentOid  string
	CanDelete canDeleted
}

type canDeleted struct {
	AgentcanDeleted bool
	Reason          string
}

func (db *UsersDb) Condition() string {
	return db.UserOid
}
func (db *UsersDb) CanDeleted() bool {
	return db.CanDelete.AgentcanDeleted
}
