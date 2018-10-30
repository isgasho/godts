package main

import (
	"fmt"
	"github.com/DaigangLi/godts/cache"
	"github.com/DaigangLi/godts/conf"
	"github.com/DaigangLi/godts/db"
	"github.com/DaigangLi/godts/zk"
	"log"
)

func init() {
}

func main() {

	prop := &conf.Prop{}

	zkServers := prop.GetString("zookeeper.connect")
	//zksessionTimeout := prop.GetInt("zookeeper.connection.timeout.ms", 5000)

	//zkClient := &zk.ZkClient{zkServers, zksessionTimeout}
	//zkClient.Connect()

	// zookeeper配置
	zkConfig := &zk.ZookeeperConfig{
		Servers:    []string{zkServers},
		RootPath:   "/godts",
		MasterPath: "/master",
	}

	// main goroutine 和 选举goroutine之间通信的channel，同于返回选角结果
	isMasterChan := make(chan bool)

	var isMaster bool

	// 选举
	electionManager := zk.NewElectionManager(zkConfig, isMasterChan)
	go electionManager.Run()

	for {
		select {
		case isMaster = <-isMasterChan:
			if isMaster {
				// 开启DTS任务
				// 1.从元数据中恢复集群状态
				// 2.重置各个数据源的同步数据位置
				// 3.开始同步数据
				log.Println("I'm master. wait for start dts service.")

				yml := &conf.Yml{}
				ymlContext := yml.GetYmlContext()

				db.InitEngine(ymlContext.Mysql)
				results, err := db.GetEngine().QueryString("show master status")
				if err == nil {
					for k, v := range results {
						fmt.Printf("%s=%d;", k, v)
					}
				}

				cache.NewCache()
				cache.Cache().SetDefault("test", ymlContext.Mysql)

				if item, ok := cache.Cache().Get("test"); ok {
					if mysqlConf, ok := item.(*conf.Mysql); ok {
						log.Printf(mysqlConf.BinlogFile)
					}
				}

				//db.Select()
				//

				//
				////mail.Send(ymlContext.Mail)
				//
				//source.StartCanal(ymlContext.Mysql)

				// 开启web
				//web.Start()

			}
		}
	}
}
