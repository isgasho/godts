package zk

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"strings"
	"time"
)

type ZkConfig struct {
	servers        []string
	sessionTimeout time.Duration
}

func NewZkConfig(servers string, sessionTimeout int) *ZkConfig {
	zkConfig := &ZkConfig{
		strings.Split(servers, ","),
		time.Millisecond * time.Duration(sessionTimeout),
	}
	return zkConfig
}

type ZkClient struct {
	zkConfig *ZkConfig
	conn     *zk.Conn
}

func (zkClient *ZkClient) Connect() {

	zkConfig := zkClient.zkConfig
	conn, event, err := zk.Connect(zkConfig.servers, zkConfig.sessionTimeout)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	//// get
	//data, state, err := conn.Get("/")
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//
	//log.Println("%s", data)
	//log.Println("%s", ZkStateStringFormat(state))

	// 等待连接成功
	for {
		isConnected := false
		select {
		case zkEvent := <-event:
			if zkEvent.State == zk.StateConnected {
				isConnected = true
				log.Println("connect to zookeeper server success!")
			}
		case _ = <-time.After(zkConfig.sessionTimeout):
			log.Println("connect to zookeeper server timeout!")
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
