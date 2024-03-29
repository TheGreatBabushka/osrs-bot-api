package bot

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type Bot struct {
	ID       string   `json:"id"`
	Email    string   `json:"email"`    // dreambot username (email)
	Username string   `json:"username"` // osrs username
	Script   string   `json:"script"`
	Params   []string `json:"params"`
	Status   string   `json:"status"`
	PID      int      `json:"pid"` // process id of the dreambot client, 0 if not running. set by the server when bot is started
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
	if b.PID == 0 {
		return false
	}

	cmd := exec.Command("cmd", "/C", "tasklist", "/FI", fmt.Sprintf("PID eq %d", b.PID))
	out, err := cmd.Output()
	if err != nil {
		fmt.Print(err, "\n")
		return false
	}

	output := string(out)
	if output == "" || strings.Contains(output, "No tasks are running which match the specified criteria") {
		return false
	}

	return true
}

func (b *Bot) startDreamBotClient() {
	log.Println("Starting DreamBot client for Bot: " + b.Email + " with script: " + b.Script + "")

	client_path := "C:\\Users\\Administrator\\DreamBot\\BotData\\client.jar"

	var clientParams = []string{"-jar", client_path, "-account", b.Email, "-script", b.Script, "-world", "f2p", "-covert", "-fresh"}

	// chech for bot/script specific params
	if b.Params != nil && len(b.Params) > 0 {
		log.Println("Found bot specific params: " + fmt.Sprint(b.Params))

		// if params doesnt start with -params, insert it at the beginning
		if !strings.HasPrefix(b.Params[0], "-params") {
			clientParams = append(clientParams, "-params")
		}

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
