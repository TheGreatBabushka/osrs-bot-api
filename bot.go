package main

import (
	"fmt"
	"log"
	"os/exec"
)

type bot struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`     // dreambot username (email)
	Username string   `json:"username"` // osrs username
	Script   string   `json:"script"`
	Params   []string `json:"params"`
	Status   string   `json:"status"`
	PID      int      `json:"pid"`
}

type account_levels struct {
	Attack      int `json:"attack"`
	Strength    int `json:"strength"`
	Defence     int `json:"defence"`
	Ranged      int `json:"ranged"`
	Magic       int `json:"magic"`
	Prayer      int `json:"prayer"`
	Runecraft   int `json:"runecraft"`
	Hitpoints   int `json:"hitpoints"`
	Agility     int `json:"agility"`
	Herblore    int `json:"herblore"`
	Thieving    int `json:"thieving"`
	Crafting    int `json:"crafting"`
	Fletching   int `json:"fletching"`
	Slayer      int `json:"slayer"`
	Hunter      int `json:"hunter"`
	Mining      int `json:"mining"`
	Smithing    int `json:"smithing"`
	Fishing     int `json:"fishing"`
	Cooking     int `json:"cooking"`
	Firemaking  int `json:"firemaking"`
	Woodcutting int `json:"woodcutting"`
	Farming     int `json:"farming"`
}

type heartbeat struct {
	Username string         `json:"username"` // dreambot username / osrs login email
	Status   string         `json:"status"`
	Levels   account_levels `json:"levels"`
}

func (b *bot) startBot() {
	if b.Status == "Started" {
		return
	}

	b.Status = "Started"
	b._startDreamBotClient()
}

func (b *bot) stopBot() {
	if b.Status == "Stopped" {
		return
	}

	b.Status = "Stopped"

	log.Println("Stopping DreamBot client for bot=" + b.Name + " (currently running script: " + b.Script + ")")
	cmd := exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprint(b.PID))
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Client stopped with pid %d", cmd.Process.Pid)
}

func (b *bot) restartBot() {
	b.Status = "Restarting"
}

func (b *bot) _startDreamBotClient() {
	log.Println("Starting DreamBot client for Bot: " + b.Name + " with script: " + b.Script + "")

	client_path := "C:\\Users\\Administrator\\DreamBot\\BotData\\client.jar"

	var clientParams = []string{"-jar", client_path, "-account", b.Name, "-script", b.Script, "-world", "f2p", "-covert", "-fresh"}

	// chech for bot/script specific params
	if b.Params != nil {
		log.Println("Found bot specific params: " + fmt.Sprint(b.Params))
		clientParams = append(clientParams, b.Params...)
	}

	cmd := exec.Command("java", clientParams...)
	log.Println("Starting DreamBot client with command: java " + fmt.Sprint(clientParams))
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Client started with pid %d", cmd.Process.Pid)
	b.PID = cmd.Process.Pid
}
