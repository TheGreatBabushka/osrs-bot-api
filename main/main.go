package main

import (
	"bot-api/api"
	s "bot-api/server"
)

// var db d.Database
// var sql_db *sql.DB

func main() {
	server := s.Server{}
	server.Start()

	api.Start(&server)
}
