package main

import (
	"bot-api/api"
	s "bot-api/server"
)

func main() {
	server := s.Server{}
	api.Start(&server)
}
