package main

import (
	"github.com/DaigangLi/godts/src/dts/prop"
	"github.com/DaigangLi/godts/src/dts/zk"
	"time"
)

func main() {

	zks := prop.GetString("zookeeper.connect")
	sessionTimeout := prop.GetInt("zookeeper.connection.timeout.ms", 5000)
	zk.Connect([]string{zks}, time.Millisecond*time.Duration(sessionTimeout))
}
