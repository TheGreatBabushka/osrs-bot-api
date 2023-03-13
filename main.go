package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var bots = []bot{}
var latestHeartbeats = map[string]heartbeat{}

func main() {
	router := gin.Default()

	router.GET("/bots", getBots)
	router.GET("/bots/:id", getBotByID)

	router.POST("/heartbeat", handleHeartbeat)
	router.PUT("/bots", putBots)

	router.DELETE("/bots/:id", deleteBot)

	router.Run("192.168.1.156:8080")
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

	var hb heartbeat
	if err := c.BindJSON(&hb); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Printf("Error parsing heartbeat: %s\n", c.Request.Body)
		return
	}

	for i, b := range bots {
		if b.ID == hb.ID {
			bots[i].Status = hb.Status

			hb_changed := false
			if latestHeartbeats[hb.ID].Status != hb.Status {
				hb_changed = true
			}
			latestHeartbeats[hb.ID] = hb

			if hb_changed {
				fmt.Println("Bot " + hb.ID + " status has changed to: " + hb.Status)
			}

			c.IndentedJSON(http.StatusOK, gin.H{"message": "heartbeat received"})
			return
		}
	}

	fmt.Println("Heartbeat received from unknown bot with id: " + hb.ID)
	bots = append(bots, bot{ID: hb.ID, Status: hb.Status})
}

// REST API

// GET
// /bots
// /bots/:id

// POST
// /bots
