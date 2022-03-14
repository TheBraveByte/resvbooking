package driver

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"time"
)

const (
	maxDbOpenConn = 10
	maxDbLifeTime = 5 * time.Minute
	maxDbIdConn   = 5
	driverName    = "pgx"
)

type DB struct {
	PSQL *sql.DB
}

var dbConn = &DB{}

func ConnectSqlDb(dsn string) (*DB, error) {
	db, err := NewDatabase(dsn)

	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(maxDbOpenConn)
	db.SetMaxIdleConns(maxDbIdConn)
	db.SetConnMaxLifetime(maxDbLifeTime)

	dbConn.PSQL = db
	return dbConn, nil
}

func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open(driverName, dsn)

	if err != nil {
		log.Fatal(fmt.Sprintln("Error connecting to database"))
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("Database Pinged Successfully")

	return db, nil
}
