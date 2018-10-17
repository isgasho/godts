package main

import (
	"fmt"
	"github.com/DaigangLi/godts/src/dts/prop"
	"github.com/DaigangLi/godts/src/dts/zk"
)

func main() {

	prop := &prop.Prop{}
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
				fmt.Println("I'm master. wait for start dts service.")
			}
		}
	}
}
