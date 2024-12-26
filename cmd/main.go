package main

import (
	"log"
	"rest_api_pks/internal/server"
)

func main() {
	server := server.NewServer()
	log.Fatal(server.Start(":8080"))
}
