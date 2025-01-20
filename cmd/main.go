package main

import (
	"log"

	"github.com/diogocardoso/go/lab_1/configs"
	"github.com/diogocardoso/go/lab_1/infrastructure"
)

func main() {
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %s", err.Error())
	}

	server := infrastructure.NewServer(config)
	server.Start()
}
