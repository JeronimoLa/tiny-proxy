package db

import (
	"fmt"
	"database/sql"
	"sync"
	"time"
	"log"
	"net/http"
	"net/http/httputil"
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

func StartRequestDB() (*dbConnection){
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

// grab urls and group them in transit to feed into AI model or
// index my database and feed the url/domains into the AI model 

// func 

func printLogger(data *ReqLogger){
	// fmt.Println(data.timestamp + " | " + data.source_ip +" | " + data.method + " | " + data.url + " | " + data.url_domain + " | " + data.user_agent)
	fmt.Println(data.method + " | " + data.url + " | " + data.url_domain + " | " + data.user_agent)
}

func RequestLogger(connection *dbConnection, req *http.Request) {
	raw_request, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Fatal(err)
	}

	dataDump := &ReqLogger{
		timestamp: time.Now().UTC().Format(time.RFC1123),
		source_ip: req.RemoteAddr,
		method: req.Method,
		url: req.URL.String(),
		url_domain: req.Host,
		user_agent: req.UserAgent(), 
		raw_data: string(raw_request),	
	}
	insertRequestDataToDB(connection, dataDump)
	printLogger(dataDump)
}