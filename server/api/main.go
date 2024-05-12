package main

import (
	"github.com/tom-draper/api-analytics/server/api/lib/log"
	"github.com/tom-draper/api-analytics/server/api/lib/routes"
	"github.com/tom-draper/api-analytics/server/database"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	log.LogToFile("Starting api...")

	log.LogToFile("Running Migrations...")

	tablesError := database.CreateTables()

	if tablesError != nil {
		log.LogToFile("Error creating tables: " + tablesError.Error())
	}

	gin.SetMode(gin.ReleaseMode)
	app := gin.New()

	r := app.Group("/api")

	r.Use(cors.Default())

	routes.RegisterRouter(r)

	app.Run(":3000")
}
