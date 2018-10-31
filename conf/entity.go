package conf

import (
	"fmt"
)

type YmlContext struct {
	Mail  *Mail  `yaml:"mail"`
	Mysql *Mysql `yaml:"mysql"`
}

type Mail struct {
	Host     string
	Port     int
	UserName string `yaml:"user-name"`
	Password string
	SslPort  int `yaml:"ssl-port"`
}

type DBConf interface {
	ToConnectStr() string
}

type Mysql struct {
	ServerId   uint32 `yaml:"server-id"`
	Flavor     string
	Host       string
	Port       uint16
	User       string
	Password   string
	Database   string
	BinlogFile string `yaml:"binlog-file"`
}

func (mysqlConf *Mysql) ToConnectStr() string {
	port := fmt.Sprint(mysqlConf.Port)
	return mysqlConf.User + ":" + mysqlConf.Password + "@tcp(" + mysqlConf.Host + ":" + port + ")/" + mysqlConf.Database
}
