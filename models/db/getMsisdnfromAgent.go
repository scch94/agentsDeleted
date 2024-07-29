package modeldb

import "database/sql"

type MsisdnDbSql struct {
	Msisdn    sql.NullString
	MsisdnOid sql.NullString
	AgentOid  sql.NullString
}

func (s *MsisdnDbSql) ConvertMsisdn() MsisdnDb {
	return MsisdnDb{
		Msisdn:    s.Msisdn.String,
		MsisdnOid: s.MsisdnOid.String,
		AgentOid:  s.AgentOid.String,
	}
}

type MsisdnDb struct {
	Msisdn    string
	MsisdnOid string
	AgentOid  string
}

func (db *MsisdnDb) Condition() string {
	return db.MsisdnOid
}

func (db *MsisdnDb) CanDeleted() bool {
	return true
}
