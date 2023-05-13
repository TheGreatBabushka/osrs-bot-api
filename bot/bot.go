package bot

import (
	"fmt"
	"log"
	"os/exec"
)

type Bot struct {
	ID       string   `json:"id"`
	Email    string   `json:"name"`     // dreambot username (email)
	Username string   `json:"username"` // osrs username
	Script   string   `json:"script"`
	Params   []string `json:"params"`
	Status   string   `json:"status"`
	PID      int      `json:"pid"`
}

func (b *Bot) Start() {
	if b.Status == "Started" {
		return
	}

	b.Status = "Started"
	b.startDreamBotClient()
}

func (b *Bot) Stop() {
	if b.Status == "Stopped" {
		return
	}

	b.Status = "Stopped"

	log.Println("Stopping DreamBot client for bot=" + b.Email + " (currently running script: " + b.Script + ")")
	cmd := exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprint(b.PID))
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Client stopped with pid %d", cmd.Process.Pid)
}

func (b *Bot) IsRunning() bool {
	// TODO - this seems to be broken, use a different method

	command := "tasklist /FI \"PID eq " + fmt.Sprint(b.PID) + "\""
	cmd := exec.Command("cmd", "/C", command)
	fmt.Printf("Running command: %s\n", command)

	out, err := cmd.Output()
	if err != nil {
		fmt.Print(err, "\n")
		return false
	}

	if string(out) == "" {
		return false
	}

	if b.PID == 0 {
		return false
	}

	return true
}

func (b *Bot) startDreamBotClient() {
	log.Println("Starting DreamBot client for Bot: " + b.Email + " with script: " + b.Script + "")

	client_path := "C:\\Users\\Administrator\\DreamBot\\BotData\\client.jar"

	var clientParams = []string{"-jar", client_path, "-account", b.Email, "-script", b.Script, "-world", "f2p", "-covert", "-fresh"}

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
