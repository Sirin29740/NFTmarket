package main

import (
	"NFTmarket/api"
	"NFTmarket/internal/database"
)

func main() {
	database.InitDB()
	api.Router()
}
