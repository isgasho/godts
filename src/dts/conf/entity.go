package conf

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

type Mysql struct {
	ServerId   uint32 `yaml:"server-id"`
	Flavor     string
	Host       string
	Port       uint16
	User       string
	Password   string
	BinlogFile string `yaml:"binlog-file"`
}
