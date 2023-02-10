package db

import (
	"database/sql"
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
	conn, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/signavault")

	conn.SetConnMaxLifetime(time.Minute * 3)
	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(10)

	if err != nil {
		panic(err.Error())
	}
	return conn
}
