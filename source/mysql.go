package source

import (
	"context"
	"fmt"
	"github.com/DaigangLi/godts/conf"
	"github.com/siddontang/go-mysql/canal"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartReplication(mysqlConf *conf.Mysql) {

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

	var binlogPos uint32 = 0

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

		ev.Header.Dump(os.Stdout)
		if rowsEvent, ok := ev.Event.(*replication.RowsEvent); ok {
			fmt.Fprintf(os.Stdout, "Table Name:%q\n", rowsEvent.Table.Table)
			fmt.Fprintf(os.Stdout, "Schema:%q\n", rowsEvent.Table.Schema)
			for _, rows := range rowsEvent.Rows {
				fmt.Fprintf(os.Stdout, "--\n")
				for j, d := range rows {
					if _, ok := d.([]byte); ok {
						fmt.Fprintf(os.Stdout, "%d:%q\n", j, d)
					} else {
						fmt.Fprintf(os.Stdout, "%d:%#v\n", j, d)
					}
				}
			}
		}

		//ev.Dump(os.Stdout)
	}
}

func StartCanal(mysqlConf *conf.Mysql) {
	cfg := canal.NewDefaultConfig()
	cfg.Addr = fmt.Sprintf("%s:%d", mysqlConf.Host, mysqlConf.Port)
	cfg.User = mysqlConf.User
	cfg.Password = mysqlConf.Password

	c, err := canal.NewCanal(cfg)
	if err != nil {
		fmt.Printf("create canal err %v", err)
		os.Exit(1)
	}

	c.SetEventHandler(&handler{})

	startPos := mysql.Position{
		Name: mysqlConf.BinlogFile,
		Pos:  0,
	}

	go func() {
		err = c.RunFrom(startPos)
		if err != nil {
			fmt.Printf("start canal err %v", err)
		}
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		os.Kill,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	<-sc

	c.Close()
}

type handler struct {
	canal.DummyEventHandler
}

func (h *handler) OnRow(e *canal.RowsEvent) error {
	fmt.Printf("%v\n", e)

	return nil
}

func (h *handler) String() string {
	return "TestHandler"
}
