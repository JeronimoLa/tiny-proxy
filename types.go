
package main


type ReqLogger struct {
	timestamp	string
	source_ip	string
	method		string
	url 		string
	user_agent	string	
	raw_data	[]byte
}