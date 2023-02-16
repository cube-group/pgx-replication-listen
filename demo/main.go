package main

import (
	"fmt"
	"github.com/cube-group/pgx-replication/core"
	"github.com/jackc/pgx"
	"io/ioutil"
	"log"
	"strconv"
)

var syncer *core.ReplicationSyncer

func main() {
	syncer = core.NewReplicationSyncer(
		//pgx.ConnConfig{
		//	Host:     "127.0.0.1",
		//	Port:     4321,
		//	Database: "default",
		//	User:     "postgres",
		//	Password: "",
		//},
		pgx.ConnConfig{
			Host:     "192.168.4.157",
			Port:     30433,
			Database: "web",
			User:     "postgres",
			Password: "default",
		},
		handle,
	)

	var wal int64
	bytes, err := ioutil.ReadFile(".wal")
	if err == nil {
		wal, _ = strconv.ParseInt(string(bytes), 0, 64)
	}
	log.Fatalf("sync err: %v", syncer.Debug().Start("test_pg_2_ck", []string{"test"}, uint64(wal)))
}

func handle(msg core.ReplicationMessage) {
	if msg.Wal > 0 { // current wal for startLSN
		ioutil.WriteFile(".wal", []byte(fmt.Sprintf("%v", msg.Wal)), 0777)
	}
	if msg.Namespace != "" {
		log.Printf("[%v.%v] (%v) %+v", msg.Namespace, msg.Table, msg.MessageType, msg.Body)
	}
}
