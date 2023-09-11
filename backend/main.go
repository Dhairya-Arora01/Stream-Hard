package main

import (
	"github.com/Dhairya-Arora01/StreamHard/server"
	"github.com/Dhairya-Arora01/StreamHard/server/db"
)

func main() {
	db.InitDB()
	server.CreateRoutes("8000")

}
