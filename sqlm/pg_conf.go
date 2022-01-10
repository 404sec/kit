package sqlm

import "gorm.io/gorm"

type PGDB struct {
	Client *gorm.DB
	Cfg    *PGConfig
}

func (c PGDB) GetConfig() interface{} {
	return c.Cfg
}

type PGConfig struct {
	Clf DBLoginConfig
	Con DbConfigConn
}
type DBLoginConfig struct {
	Host        string
	TimeZone    string
	Port        string
	User        string
	Password    string
	Dbname      string
	Sslcert     string
	Sslkey      string
	Sslrootcert string
	Sslmode     string
}
type DbConfigConn struct {
	MaxIdleConns       int
	SetMaxOpenConns    int
	LongPollerInterval int64 //轮询时间
}
