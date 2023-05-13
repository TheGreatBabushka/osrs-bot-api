package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	b "bot-api/bot"
	s "bot-api/server"
)

var server s.Server

func Start(s *s.Server) {
	server = *s

	router := gin.Default()

	router.GET("/bots", getBots)
	router.PUT("/bots", startBot)
	router.GET("/bots/:id", getBotByID)
	router.DELETE("/bots/:id", deleteBot)

	router.POST("/heartbeat", handleHeartbeat)

	router.GET("/accounts", getAccounts)
	router.GET("/accounts/:id", getAccountByID)

	router.GET("/levels/:id", getLevelsByID)

	router.Run("192.168.1.156:8080")
}

func getBots(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, server.Bots)
}

func startBot(c *gin.Context) {
	var newBot b.Bot
	if err := c.BindJSON(&newBot); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ensure ID is not empty and is unique
	if newBot.ID == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "ID is empty"})
		return
	}

	for _, bot := range server.Bots {
		if bot.ID == newBot.ID {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "ID is not unique"})
			return
		}
	}

	newBot.Status = "Stopped"
	newBot.Start()

	server.Bots = append(server.Bots, newBot)
	c.IndentedJSON(http.StatusCreated, newBot)
}

func getBotByID(c *gin.Context) {
	id := c.Param("id")

	for _, b := range server.Bots {
		if b.ID == id {
			c.IndentedJSON(http.StatusOK, b)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func deleteBot(c *gin.Context) {
	id := c.Param("id")

	for i, b := range server.Bots {
		if b.ID == id {
			b.Stop()
			server.Bots = append(server.Bots[:i], server.Bots[i+1:]...)
			c.IndentedJSON(http.StatusOK, gin.H{"message": "bot deleted"})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "bot not found"})
}

func handleHeartbeat(c *gin.Context) {

	// parse heartbeat
	var hb b.Heartbeat
	if err := c.BindJSON(&hb); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Printf("Error parsing heartbeat: %s\n", c.Request.Body)
		return
	}

	// TODO - move logic to server
	// check if bot is known
	for i, b := range server.Bots {
		if b.Email == hb.Email {
			server.Bots[i].Status = hb.Status

			hb_changed := false
			if server.LatestHeartbeats[hb.Email].Status != hb.Status {
				hb_changed = true
			}

			server.LatestHeartbeats[hb.Email] = hb

			if hb_changed {
				fmt.Println("Bot " + hb.Email + " status has changed to: " + hb.Status)
			}

			account, err := server.DB.GetAccountByEmail(hb.Email)
			if err != nil {
				if err == sql.ErrNoRows {
					fmt.Println("Account not found for username: " + hb.Email)

				}
				fmt.Println("Error getting account for username: " + hb.Email)
				fmt.Println(err)
				return
			}

			err = server.DB.UpdateLevelsForAccount(account, hb.Levels)
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
	server.Bots = append(server.Bots, b.Bot{Email: hb.Email, Status: hb.Status})
	fmt.Println("Levels: " + fmt.Sprint(hb.Levels) + "\n")

	fmt.Println("Adding bot to database: " + hb.Email)
	server.DB.InsertAccount(hb.Email, hb.Username)
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
	id_int, err := strconv.Atoi(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a, err := server.DB.GetAccount(id_int)
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

	a, err := server.DB.GetAccount(id_int)
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
