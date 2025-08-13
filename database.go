package main

import (
	"fmt"
	"database/sql"
	"sync"
    _ "github.com/glebarez/go-sqlite"
)

var ( 
	once sync.Once
	dbConn *dbConnection
)

func getDBInstance() (*dbConnection, error) {
	var initError error
	once.Do(func() {
		db, err := sql.Open("sqlite", "requests.db")
		if err != nil {
			initError = fmt.Errorf("failed to open database: %v", err)
			return
		}
		dbConn = &dbConnection{db: db}
	})
	if initError != nil {
		return nil, initError
	}
	return dbConn, nil
}

func createRequestTable(connection *dbConnection) {
	sql_table := `CREATE TABLE IF NOT EXISTS internal_network (
		id INTEGER PRIMARY KEY,
		timestamp TEXT NOT NULL,
		source_ip TEXT NOT NULL,
		method TEXT NOT NULL,
		url TEXT NOT NULL,
		url_domain TEXT NOT NULL,
		user_agent TEXT NOT NULL,
		raw_data TEXT NOT NULL
	);`
	connection.db.Exec(sql_table)
	fmt.Println("Creating table if it doesn't exist")
}

func startRequestDB() (*dbConnection){
	conn, err := getDBInstance()
	if err != nil {
		fmt.Print(err)
	}
	createRequestTable(conn)
	return conn
}

func insertRequestDataToDB(connection *dbConnection, log *ReqLogger) {

	insert_sql := `INSERT INTO internal_network 
		(timestamp,
		source_ip,
		method,
		url,
		url_domain,
		user_agent,
		raw_data) VALUES (?, ?, ?, ?, ?, ?, ?);`
	
	result, err := connection.db.Exec(insert_sql, log.timestamp, log.source_ip, log.method, log.url, log.url_domain, log.user_agent, log.raw_data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result.LastInsertId())
}