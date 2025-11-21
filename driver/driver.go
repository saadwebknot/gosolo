package driver

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// DB ...
type DB struct {
	SQL *sql.DB
	// Mgo *mgo.database
}

// DBConn ...
var dbConn = &DB{}

// ConnectSQL ...
func ConnectSQL(host, port, uname, pass, dbname string) (*DB, error) {
	// TiDB Cloud requires TLS/SSL - add tls=true to connection string
	dbSource := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&loc=UTC&tls=true",
		uname,
		pass,
		host,
		port,
		dbname,
	)
	fmt.Printf("Connecting to database: %s@%s:%s/%s\n", uname, host, port, dbname)
	d, err := sql.Open("mysql", dbSource)
	if err != nil {
		return nil, err
	}
	// Test the connection
	if err := d.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}
	fmt.Printf("Successfully connected to database: %s\n", dbname)
	dbConn.SQL = d
	return dbConn, nil
}
