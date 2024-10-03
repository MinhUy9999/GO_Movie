package main

import (
	"my-app/config"
	_ "my-app/docs" // Import Swagger docs
	"my-app/routes"

	"github.com/gin-contrib/cors"
)

func main() {
	// Connect to the database
	config.ConnectDB()

	// Set up routes
	r := routes.SetupRoutes()
	r.Use(cors.Default())
	// Start server on port 8080
	r.Run(":8080")
}
