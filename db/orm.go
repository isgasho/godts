package db

import (
	"github.com/DaigangLi/godts/conf"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xormplus/xorm"
)

var engine *xorm.Engine

func Init(mysqlConf *conf.Mysql) {
	newEngine, err := xorm.NewEngine("mysql", mysqlConf.ToConnectStr())
	if err != nil {
		engine = newEngine
	}

	engine.ShowSQL(true)
	engine.DBMetas()
}
