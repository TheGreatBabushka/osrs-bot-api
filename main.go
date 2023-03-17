package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// list of all bots currently running and known to the server
var bots = []bot{}

// map for a bot's username to the last received heartbeat for that bot
var latestHeartbeats = map[string]heartbeat{}

var db database
var sql_db *sql.DB

func main() {
	sql_db = initDatabase()
	db = database{sql_db}

	router := gin.Default()

	router.GET("/bots", getBots)
	router.PUT("/bots", putBots)
	router.GET("/bots/:id", getBotByID)
	router.DELETE("/bots/:id", deleteBot)

	router.POST("/heartbeat", handleHeartbeat)

	router.GET("/accounts", getAccounts)
	router.GET("/accounts/:id", getAccountByID)

	router.GET("/levels/:id", getLevelsByID)

	router.Run("192.168.1.156:8080")
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

func getBots(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, bots)
}

func putBots(c *gin.Context) {
	var newBot bot
	if err := c.BindJSON(&newBot); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ensure ID is not empty and is unique
	if newBot.ID == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "ID is empty"})
		return
	}

	for _, bot := range bots {
		if bot.ID == newBot.ID {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "ID is not unique"})
			return
		}
	}

	newBot.Status = "Stopped"
	newBot.startBot()

	bots = append(bots, newBot)
	c.IndentedJSON(http.StatusCreated, newBot)
}

func getBotByID(c *gin.Context) {
	id := c.Param("id")

	for _, b := range bots {
		if b.ID == id {
			c.IndentedJSON(http.StatusOK, b)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func deleteBot(c *gin.Context) {
	id := c.Param("id")

	for i, b := range bots {
		if b.ID == id {
			b.stopBot()
			bots = append(bots[:i], bots[i+1:]...)
			c.IndentedJSON(http.StatusOK, gin.H{"message": "bot deleted"})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "bot not found"})
}

func handleHeartbeat(c *gin.Context) {

	// parse heartbeat
	var hb heartbeat
	if err := c.BindJSON(&hb); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Printf("Error parsing heartbeat: %s\n", c.Request.Body)
		return
	}

	// check if bot is known
	for i, b := range bots {
		if b.Email == hb.Email {
			bots[i].Status = hb.Status

			hb_changed := false
			if latestHeartbeats[hb.Email].Status != hb.Status {
				hb_changed = true
			}

			latestHeartbeats[hb.Email] = hb

			if hb_changed {
				fmt.Println("Bot " + hb.Email + " status has changed to: " + hb.Status)
			}

			account, err := db.getAccountByEmail(hb.Email)
			if err != nil {
				fmt.Println("Error getting account for username: " + hb.Email)
				fmt.Println(err)
				return
			}

			err = db.updateLevelsForAccount(account, hb.Levels)
			if err != nil {
				fmt.Println("Error updating levels for account: " + account.Username)
				fmt.Println(err)
				return
			}

			c.IndentedJSON(http.StatusOK, gin.H{"message": "heartbeat received"})
			return
		}
	}

	// bot is not known, add it to the list of known bots
	fmt.Println("Heartbeat received from unknown bot with username: " + hb.Email)
	bots = append(bots, bot{Email: hb.Email, Status: hb.Status})

	fmt.Println("Levels: " + fmt.Sprint(hb.Levels) + "\n")
}

func getAccounts(c *gin.Context) {
	accounts, err := db.getAllAccounts()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, accounts)
}

func getAccountByID(c *gin.Context) {

	id := c.Param("id")
	id_int, err := strconv.Atoi(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a, err := db.getAccount(id_int)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, a)
}

func getLevelsByID(c *gin.Context) {

	id := c.Param("id")
	id_int, err := strconv.Atoi(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a, err := db.getAccount(id_int)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	levels, err := db.getLevelsForAccount(a.ID)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, levels)
}

// REST API

// GET
// /bots
// /bots/:id

// POST
// /bots
