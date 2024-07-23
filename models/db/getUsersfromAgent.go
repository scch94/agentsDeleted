package modeldb

type UsersDb struct {
	UserOid   string
	UserId    string
	AgentOid  string
	CanDelete bool
}

func (db *UsersDb) Condition() string {
	return db.UserOid
}
