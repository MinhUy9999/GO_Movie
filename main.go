package main

import (
	"log"
	"my-app/config"
	_ "my-app/docs" // Import Swagger docs
	"my-app/routes"

	"github.com/gin-contrib/cors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Connect to the database
	config.ConnectDB()

	// Connect to Redis
	err := config.ConnectRedis()
	if err != nil {
		log.Printf("Cảnh báo: Không thể kết nối đến Redis: %v", err)
		// Tiếp tục với ứng dụng, nhưng không có chức năng Redis
	} else {
		defer config.RedisClient.Close()
		log.Println("Kết nối thành công đến Redis")
	}

	// Set up routes
	r := routes.SetupRoutes()

	// Add Swagger documentation route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Use CORS middleware
	r.Use(cors.Default())

	// Start server on port 8080
	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
