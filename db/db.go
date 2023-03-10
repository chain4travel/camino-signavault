package db

import (
	"database/sql"
	"github.com/chain4travel/camino-signavault/util"
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"time"
)

var lock = &sync.Mutex{}

type Db struct {
	*sql.DB
}

var dbInstance *Db

func GetInstance() *Db {
	if dbInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if dbInstance == nil {
			connection := initConnection()
			dbInstance = &Db{connection}
		}
	}
	return dbInstance
}

func initConnection() *sql.DB {
	config := util.GetInstance()
	conn, err := sql.Open(config.Database.Type, config.Database.Dsn)
	conn.SetConnMaxLifetime(time.Minute * 3)
	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(10)

	if err != nil {
		panic(err.Error())
	}
	return conn
}
