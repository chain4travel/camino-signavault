package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr/v2"
	"sync"
	"time"
)

var lock = &sync.Mutex{}

type Db struct {
	*dbr.Connection
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

func initConnection() *dbr.Connection {
	conn, err := dbr.Open("mysql", "root:password@tcp(localhost)/signavault", nil)
	conn.SetConnMaxLifetime(time.Minute * 3)
	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(10)

	if err != nil {
		panic(err.Error())
	}
	return conn
}
