package db

import (
	"github.com/DaigangLi/godts/conf"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"sync"
)

var engine *xorm.Engine
var once sync.Once

func InitEngine(mysqlConf *conf.Mysql) {

	once.Do(func() {
		ormEngine, err := xorm.NewEngine("mysql", mysqlConf.ToConnectStr())
		if err == nil {
			engine = ormEngine
			engine.ShowSQL(true)
		}
	})
}

func GetEngine() *xorm.Engine {
	return engine
}
