package main

import (
	"file-taiko-wiki/src/server"
	"file-taiko-wiki/src/util"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	util.LoadFlag()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	app := server.CreateServer()
	app.Run(os.Getenv("PORT"))
}
