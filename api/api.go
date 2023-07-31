package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	b "bot-api/bot"
	s "bot-api/server"
)

var server s.Server

type StartBotCommand struct {
	ID     string   `json:"id"`
	Script string   `json:"script"`
	Params []string `json:"params"`
}

func Start(s *s.Server) {
	server = *s

	router := gin.Default()

	router.GET("/bots/active", getActiveBots)
	router.GET("/bots/inactive", getInactiveBots)

	router.GET("/bots/activity", getBotActivity)
	router.GET("/bots/activity/:id", getBotActivityByID)
	// router.GET("/bots/heartbeat", getBotHeartbeat)

	router.POST("/bots", startBot)
	router.GET("/bots/:id", getBotByID)
	router.DELETE("/bots/:id", deleteBot)

	router.POST("/heartbeat", handleHeartbeat)

	router.GET("/accounts", getAccounts)
	router.GET("/accounts/:id", getAccountByID)

	router.GET("/levels/:id", getLevelsByID)

	server.Start()

	// needs to be the last line in the function
	router.Run("localhost:8080")
}

// return info on all bots currently running on the server
func getActiveBots(c *gin.Context) {
	bots, err := server.DB.GetActiveBots()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, bots)
}

func getInactiveBots(c *gin.Context) {
	bots, err := server.DB.GetInactiveBots()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, bots)
}

// starts a new dreambot client with the given parameters
func startBot(c *gin.Context) {
	var startCmd StartBotCommand
	if err := c.BindJSON(&startCmd); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if startCmd.ID == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "ID is empty"})
		return
	}

	acc, err := server.DB.GetAccount(startCmd.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Account not found for ID: " + startCmd.ID})
			return
		}
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bots, err := server.DB.GetActiveBots()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, b := range bots {
		if b.ID == startCmd.ID {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "bot with given ID is already running"})
			return
		}
	}

	var newBot b.Bot
	newBot.ID = fmt.Sprint(acc.ID)
	newBot.Username = acc.Username
	newBot.Email = acc.Email
	newBot.Status = "Stopped"
	newBot.Script = startCmd.Script
	newBot.Params = startCmd.Params
	newBot.Start()

	server.DB.UpdateActivity(acc.ID, startCmd.Script+" "+fmt.Sprint(startCmd.Params), newBot.PID)

	c.IndentedJSON(http.StatusCreated, startCmd)
}

// return info on a specific bot
func getBotByID(c *gin.Context) {
	id := c.Param("id")

	bots, err := server.DB.GetActiveBots()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

	if stopped := server.StopBot(id); stopped {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "bot stopped"})
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "bot not found"})
}

func handleHeartbeat(c *gin.Context) {
	// parse heartbeat
	var hb s.Heartbeat
	if err := c.BindJSON(&hb); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Printf("Error parsing heartbeat: %s\n", c.Request.Body)
		return
	}

	server.HandleHeartbeat(hb)
}

func getAccounts(c *gin.Context) {
	accounts, err := server.DB.GetAccounts()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, accounts)
}

func getAccountByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "ID is empty"})
		return
	}

	a, err := server.DB.GetAccount(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, a)
}

func getLevelsByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "ID is empty"})
		return
	}

	a, err := server.DB.GetAccount(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	levels, err := server.DB.GetLevelsForAccount(a.ID)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, levels)
}

func getBotActivity(c *gin.Context) {
	activity, err := server.DB.GetBotActivity()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, activity)
}

func getBotActivityByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "ID is empty"})
		return
	}

	activity, err := server.DB.GetBotActivityByID(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, activity)
}

// func getBotHeartbeat(c *gin.Context) {
// 	heartbeats, err := server.DB.GetHeartbeats()
// 	if err != nil {
// 		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.IndentedJSON(http.StatusOK, heartbeats)
// }
