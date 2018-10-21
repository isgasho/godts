package source

import (
	"context"
	"github.com/DaigangLi/godts/src/dts/conf"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"
	"os"
	"time"
)

func StartSync(mysqlConf *conf.Mysql) {
	// Create a binlog syncer with a unique server id, the server id must be different from other MySQL's.
	// flavor is mysql or mariadb
	cfg := replication.BinlogSyncerConfig{
		ServerID: mysqlConf.ServerId,
		Flavor:   mysqlConf.Flavor,
		Host:     mysqlConf.Host,
		Port:     mysqlConf.Port,
		User:     mysqlConf.User,
		Password: mysqlConf.Password,
	}
	syncer := replication.NewBinlogSyncer(cfg)

	binlogPos := 0
	// Start sync with specified binlog file and position
	streamer, _ := syncer.StartSync(mysql.Position{mysqlConf.BinlogFile, binlogPos})

	// or you can start a gtid replication like
	// streamer, _ := syncer.StartSyncGTID(gtidSet)
	// the mysql GTID set likes this "de278ad0-2106-11e4-9f8e-6edd0ca20947:1-2"
	// the mariadb GTID set likes this "0-1-100"

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		ev, err := streamer.GetEvent(ctx)
		cancel()

		if err == context.DeadlineExceeded {
			// meet timeout
			continue
		}

		ev.Dump(os.Stdout)
	}
}
