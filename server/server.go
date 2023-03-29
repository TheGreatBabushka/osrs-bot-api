package server

import (
	b "bot-api/bot"
	d "bot-api/db"
	"database/sql"
	"fmt"
	"time"
)

type Server struct {
	// interface to the database
	DB *d.Database

	// list of all bots currently running and known to the server
	Bots []b.Bot

	// map for a bot's username to the last received heartbeat for that bot
	LatestHeartbeats map[string]b.Heartbeat

	// stores whether the server should be running or not
	isRunning bool
}

func (s *Server) Start() {
	// initialize database
	s.DB = &d.Database{Driver: initDatabase()}

	go s.run()
}

func (s *Server) Stop() {
	s.isRunning = false
}

func (s *Server) IsRunning() bool {
	return s.isRunning
}

func (s *Server) run() {
	s.isRunning = true

	for s.isRunning {
		// fmt.Println("Server is running...")
		time.Sleep(5 * time.Second)
	}

	fmt.Println("Server has stopped.")
}

func initDatabase() *sql.DB {
	sql_db, err := sql.Open("mysql", "admin:FredLongBottoms2$@/osrs-bots")
	if err != nil {
		panic(err.Error())
	}

	err = sql_db.Ping()
	if err != nil {
		panic(err.Error())
	}

	return sql_db
}
