package sqlm

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/404sec/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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

func (c *PGDB) GetConfigEmpty() interface{} {
	return &PGConfig{}
}

func (c *PGDB) NewConnect(ctx context.Context) error {

	var err error
	d := dsn(c.Cfg.Clf)
	c.Close(ctx) //关闭旧连接
	c.Client, err = gorm.Open(postgres.Open(d), &gorm.Config{Logger: log.GetGLogger()})
	if err != nil {
		log.Errorw(ctx, "db Client err", err.Error())
		return err
	}
	sqlDB, err := c.Client.DB()
	if err != nil {
		log.Errorw(ctx, "db Client err", err.Error())
		return err
	}
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(20)
	if c.Cfg.Con.MaxIdleConns > 10 { //设置空闲连接的数量
		sqlDB.SetMaxIdleConns(c.Cfg.Con.MaxIdleConns)
	}
	if c.Cfg.Con.SetMaxOpenConns > 10 { //SetMaxOpenConns
		sqlDB.SetMaxOpenConns(c.Cfg.Con.SetMaxOpenConns)
	}

	return nil
}

func (c *PGDB) Close(ctx context.Context) {
	if c.Client != nil {
		sqlDB, err := c.Client.DB()
		if err != nil {
			log.Errorw(ctx, "db Client err", err.Error())
			return
		}
		err = sqlDB.Close()
		if err != nil {
			log.Errorw(ctx, "db close err", err.Error())
			return
		}
	}

}
func (c *PGDB) GetClient() *gorm.DB {
	return c.Client
}
