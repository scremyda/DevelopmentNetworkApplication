package main

import (
	"backened/internal/api"
	"log"
)

func main() {
	log.Println("Server started")
	api.StartServer()
	log.Println("Server shutdown")
}
