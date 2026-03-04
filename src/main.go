package main

import (
	"file-taiko-wiki/src/server"
	"file-taiko-wiki/src/util"
)

func init() {
	util.LoadFlag()
}

func main() {
	app := server.CreateServer()
	app.Run(":3000")
}
