package modeldb

type MsisdnDb struct {
	Msisdn    string
	MsisdnOid string
	AgentOid  string
}

func (db *MsisdnDb) Condition() string {
	return db.MsisdnOid
}
