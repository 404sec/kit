package sqlm

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"reflect"
	"strings"

	"github.com/404sec/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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

func dsn(i interface{}) string {
	dsn := ""
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	for k := 0; k < t.NumField(); k++ {
		if v.Field(k).Interface() != nil && v.Field(k).Interface() != "" {
			dsn = dsn + fmt.Sprintf("%v=%v ", strings.ToLower(t.Field(k).Name), v.Field(k).Interface())
		}
	}
	return dsn
}

func (c *PGConfig) NewConnect(ctx context.Context) (*gorm.DB, error) {

	var err error
	d := dsn(c)
	client, err := gorm.Open(postgres.Open(d)) //, &gorm.Config{Logger: log.GetGLogger()}
	if err != nil {
		log.Error(ctx, "db Client err", zap.Error(err))
		return client, err
	}
	sqlDB, err := client.DB()
	if err != nil {
		log.Error(ctx, "db Client err", zap.Error(err))
		return client, err
	}
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(20)
	if c.Con.MaxIdleConns > 10 { //设置空闲连接的数量
		sqlDB.SetMaxIdleConns(c.Con.MaxIdleConns)
	}
	if c.Con.SetMaxOpenConns > 10 { //SetMaxOpenConns
		sqlDB.SetMaxOpenConns(c.Con.SetMaxOpenConns)
	}

	return client, nil
}
