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

//DB database tools
type DB struct {
	PSQL *sql.DB
}

var dbConn = &DB{}

//ConnectSqlDb connect to PostgreSQL database
func ConnectSqlDb(dsn string) (*DB, error) {
	db, err := NewDatabase(dsn)

	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(maxDbOpenConn)
	db.SetMaxIdleConns(maxDbIdConn)
	db.SetConnMaxLifetime(maxDbLifeTime)

	dbConn.PSQL = db
	err = TestDatabase(db)
	if err != nil {
		return nil, err
	}
	return dbConn, nil
}

//TestDatabase testing to ping the database
func TestDatabase(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		log.Println("Error Testing Database connection")
		return err
	}
	return nil
}

//NewDatabase connect to any database tools
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
