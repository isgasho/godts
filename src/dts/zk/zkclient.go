package zk

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

type ZkClient struct {
	Servers        string
	SessionTimeout int
}

func (zkClient *ZkClient) Connect() {
	sessionTimeout := time.Millisecond * time.Duration(zkClient.SessionTimeout)

	conn, event, err := zk.Connect([]string{zkClient.Servers}, sessionTimeout)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	//// get
	//data, state, err := conn.Get("/")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//fmt.Println("%s", data)
	//fmt.Println("%s", ZkStateStringFormat(state))

	// 等待连接成功
	for {
		isConnected := false
		select {
		case zkEvent := <-event:
			if zkEvent.State == zk.StateConnected {
				isConnected = true
				fmt.Println("connect to zookeeper server success!")
			}
		case _ = <-time.After(sessionTimeout):
			fmt.Println("connect to zookeeper server timeout!")
			break
		}
		if isConnected {
			break
		}
	}
}

func ZkStateString(s *zk.Stat) string {
	return fmt.Sprintf("Czxid:%d, Mzxid: %d, Ctime: %d, Mtime: %d, Version: %d, Cversion: %d, Aversion: %d, EphemeralOwner: %d, DataLength: %d, NumChildren: %d, Pzxid: %d",
		s.Czxid, s.Mzxid, s.Ctime, s.Mtime, s.Version, s.Cversion, s.Aversion, s.EphemeralOwner, s.DataLength, s.NumChildren, s.Pzxid)
}

func ZkStateStringFormat(s *zk.Stat) string {
	return fmt.Sprintf("Czxid:%d\nMzxid: %d\nCtime: %d\nMtime: %d\nVersion: %d\nCversion: %d\nAversion: %d\nEphemeralOwner: %d\nDataLength: %d\nNumChildren: %d\nPzxid: %d\n",
		s.Czxid, s.Mzxid, s.Ctime, s.Mtime, s.Version, s.Cversion, s.Aversion, s.EphemeralOwner, s.DataLength, s.NumChildren, s.Pzxid)
}
