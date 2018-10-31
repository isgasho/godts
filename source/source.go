package source

import "github.com/DaigangLi/godts/conf"

type Source interface {
	StartReplication(conf conf.DBConf)
}
