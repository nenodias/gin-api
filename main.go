package main

import (
	"github.com/nenodias/gin-api/database"
	"github.com/nenodias/gin-api/routes"
)

func main() {
	database.ConectaComBancoDeDados()
	routes.HandleRequests()
}
