package main

import (
	"GoVideoChat-Project/internal/server"
	"log"
)

func main() {
	if err := server.Run(); err != nil {
		log.Fatalln(err)
	}
}
