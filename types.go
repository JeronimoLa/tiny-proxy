
package main

import (
	"database/sql"
)

type ReqLogger struct {
	timestamp	string
	source_ip	string
	method		string
	url 		string
	url_domain	string
	user_agent	string	
	raw_data	string
}

type dbConnection struct {
	db	*sql.DB
}