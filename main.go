package main

import (
	"log"

	"github.com/Rishi-Mishra0704/go-redis/server"
)

func main() {
	server := server.NewServer(server.Config{})
	log.Fatal(server.Start())
}
