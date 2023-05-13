package server

import (
	b "bot-api/bot"
	db "bot-api/db"
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

var bots []b.Bot

type Server struct {
	// interface to the database
	DB *db.Database

	// map for a bot's username to the last received heartbeat for that bot
	LatestHeartbeats map[string]Heartbeat

	// stores whether the server should be running or not
	isRunning bool
}

type Heartbeat struct {
	Email    string    `json:"email"`    // dreambot username / osrs login email
	Status   string    `json:"status"`   // current task status description
	Username string    `json:"username"` // osrs username
	Stats    db.Levels `json:"levels"`
	PID      int       `json:"pid"`
}

// Start the server and begin bot monitoring goroutine(s)
func (s *Server) Start() {
	// initialize database
	s.DB = &db.Database{Driver: initDatabase()}
	s.LatestHeartbeats = make(map[string]Heartbeat)

	go s.run()
}

// Stop the server and any bot monitoring goroutine(s)
func (s *Server) Stop() {
	s.isRunning = false
}

// Returns whether the server is running or not - duh
func (s *Server) IsRunning() bool {
	return s.isRunning
}

// Stops a bot and remove it the server's list of bots
func (s *Server) StopBot(id string) bool {
	// TODO - remove bot from database
	for i, b := range bots {
		if b.ID == id {
			bots = append(bots[:i], bots[i+1:]...)

			b.Stop()
			return true
		}
	}

	return false
}

func (s *Server) GetBots() []b.Bot {
	return bots
}

func (s *Server) HandleHeartbeat(hb Heartbeat) error {
	// TODO - move logic to server
	// check if bot is known
	for _, b := range bots {
		if b.Email == hb.Email {
			if err := s.handleKnownHeartbeat(hb, b); err != nil {
				fmt.Println("Error handling heartbeat for known bot: " + hb.Email)
				fmt.Println(err)
				return err
			}

			return nil
		}
	}

	// bot is not known, add it to the list of known bots
	fmt.Println("Heartbeat received from unknown bot with username: " + hb.Email)
	bots = append(bots, b.Bot{Email: hb.Email, Status: hb.Status})
	fmt.Println("Levels: " + fmt.Sprint(hb.Stats) + "\n")

	fmt.Println("Adding bot to database: " + hb.Email)
	s.DB.InsertAccount(hb.Email, hb.Username)

	return nil
}

func (s *Server) handleKnownHeartbeat(hb Heartbeat, bot b.Bot) error {
	fmt.Printf("Heartbeat received from known bot with username: %s\n", hb.Email)

	bot.Status = hb.Status

	hb_changed := false
	if s.LatestHeartbeats[hb.Email].Status != hb.Status {
		hb_changed = true
	}

	s.LatestHeartbeats[hb.Email] = hb

	if hb_changed {
		fmt.Println("Bot " + hb.Email + " status has changed to: " + hb.Status)
	}

	account, err := s.DB.GetAccountByEmail(hb.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Account not found for username: " + hb.Email)

		}
		fmt.Println("Error getting account for username: " + hb.Email)
		fmt.Println(err)
		return err
	}

	err = s.DB.UpdateLevelsForAccount(account, hb.Stats)
	if err != nil {
		fmt.Println("Error updating levels for account: " + account.Username)
		fmt.Println(err)
		return err
	}

	bots, err := s.DB.GetActiveBots()
	if err != nil {
		fmt.Println("Error getting active bots from database")
		fmt.Println(err)
		return err
	}

	for _, b := range bots {
		if b.Email == hb.Email {
			fmt.Printf("Updating activity based on heartbeat for account: %s with script: %s\n", account.Username, b.Script+" "+fmt.Sprint(b.Params))
			s.DB.UpdateActivity(account.ID, b.Script+" "+fmt.Sprint(b.Params), b.PID)
			return nil
		}
	}

	fmt.Println("Bot not found in database: " + hb.Email)
	s.DB.UpdateActivity(account.ID, "Unknown (heartbeat)", hb.PID)

	return nil
}

func (s *Server) run() {
	s.isRunning = true

	for s.isRunning {
		// fmt.Println("Server is running...")
		time.Sleep(10 * time.Second)

		s.monitorActiveBots()
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

func (s *Server) monitorActiveBots() {

	bots, err := s.DB.GetActiveBots()
	if err != nil {
		fmt.Println("Error getting active bots from database")
		fmt.Println(err)
		return
	}

	// check if each bot is still running
	// if not, update the bot's stopped_at field in the database
	// if the bot is still running, update the bot's status in the database
	for _, b := range bots {

		// check if bot is still running - for now just check if the process is still running and assume
		// the server is on the same machine as the bot
		if b.PID != 0 && b.IsRunning() {
			fmt.Printf("Bot %s is still running with PID %d\n", b.Email, b.PID)
			continue
		}

		// bot is not running, update the bot's stopped_at field in the database
		id, err := strconv.Atoi(b.ID)
		if err != nil {
			fmt.Println("Error converting bot id to int: " + b.ID)
			fmt.Println(err)
			continue
		}

		fmt.Printf("Bot %s is not running, updating stopped_at field in database\n", b.Email)
		s.DB.UpdateBotStoppedAt(id)
	}
}
