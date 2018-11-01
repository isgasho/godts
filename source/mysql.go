package source

import (
	"context"
	"fmt"
	"github.com/DaigangLi/godts/conf"
	"github.com/DaigangLi/godts/db"
	"github.com/siddontang/go-mysql/canal"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type MysqlSource struct {
}

type Databases struct {
	Id       uint64
	SchemaId int
	Name     string
	Charset  string
}

type Tables struct {
	Id         uint64
	SchemaId   uint64
	DatabaseId uint64
	Name       string
	Charset    string
	Pk         string
}

type Columns struct {
	Id           uint64
	SchemaId     uint64
	TableId      uint64
	Name         string
	Charset      string
	Coltype      string
	IsSigned     uint
	EnumValues   string
	ColumnLength uint
}

type BinlogPosition struct {
	GtidSetStr      string
	File            string
	Position        uint32
	ExecutedGtidSet string
}

func InitBinlogPosition() *BinlogPosition {
	var binlogPosition *BinlogPosition
	results, err := db.GetEngine().QueryString("show master status")

	if err == nil && len(results) > 0 {
		fristRow := results[0]
		position, _ := strconv.ParseUint(fristRow["Position"], 0, 32)

		binlogPosition = &BinlogPosition{
			GtidSetStr:      "",
			File:            fristRow["File"],
			Position:        uint32(position),
			ExecutedGtidSet: fristRow["Executed_Gtid_Set"],
		}
	}

	return binlogPosition
}

var dataBasesMap map[uint64]Databases
var tableMap map[uint64]Tables
var tableColumnMap map[string]map[int]Columns

func InitMetaData() {

	var dataBases []Databases
	db.GetEngine().Find(&dataBases)

	if len(dataBases) > 0 {
		dataBasesMap = make(map[uint64]Databases)
		for _, dataBase := range dataBases {
			dataBasesMap[dataBase.Id] = dataBase
		}
	}

	var tables []Tables
	db.GetEngine().Find(&tables)

	if len(tables) > 0 {
		tableMap = make(map[uint64]Tables)
		for _, table := range tables {
			tableMap[table.Id] = table
		}
	}

	var columns []Columns
	db.GetEngine().Find(&columns)

	if len(columns) > 0 {
		var tableId uint64
		tableId = 0

		tableColumnMap = make(map[string]map[int]Columns)

		var idx int
		for _, column := range columns {

			tlbName := tableMap[column.TableId].Name
			dbName := dataBasesMap[tableMap[column.TableId].DatabaseId].Name
			combKey := dbName + "." + tlbName

			if tableId != column.TableId {
				idx = 0
				tableId = column.TableId

				tableColumnMap[combKey] = make(map[int]Columns)
			}

			tableColumnMap[combKey][idx] = column
			idx++
		}
	}
}

func (mysqlSource *MysqlSource) StartReplication(dbConf conf.DBConf) {
	if mysqlConf, ok := dbConf.(*conf.Mysql); ok {

		InitMetaData()

		binlogPosition := InitBinlogPosition()

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

		// Start sync with specified binlog file and position
		streamer, _ := syncer.StartSync(mysql.Position{binlogPosition.File, binlogPosition.Position})

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

				findKey := string(rowsEvent.Table.Schema) + "." + string(rowsEvent.Table.Table)
				for _, rows := range rowsEvent.Rows {
					fmt.Fprintf(os.Stdout, "--\n")
					for j, d := range rows {
						if _, ok := d.([]byte); ok {
							fmt.Fprintf(os.Stdout, "%s:%q\n", tableColumnMap[findKey][j].Name, d)
						} else {
							fmt.Fprintf(os.Stdout, "%s:%#v\n", tableColumnMap[findKey][j].Name, d)
						}
					}
				}
			}

			//ev.Dump(os.Stdout)
		}
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
