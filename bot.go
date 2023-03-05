package main

import (
	"fmt"
	"log"
	"os/exec"
)

type bot struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Script string   `json:"script"`
	Params []string `json:"params"`
	Status string   `json:"status"`
	PID    int      `json:"pid"`
}

type heartbeat struct {
	ID     string `json:"id"`
	Status string `json:"status"`
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
		clientParams = append(clientParams, b.Params...)
	}

	cmd := exec.Command("java", clientParams...)
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Client started with pid %d", cmd.Process.Pid)
	b.PID = cmd.Process.Pid
}
