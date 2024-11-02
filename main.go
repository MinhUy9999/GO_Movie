package main

import (
	"log"
	"my-app/config"
	_ "my-app/docs" // Import Swagger docs
	"my-app/routes"

	// Import the sockets package
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Kết nối đến cơ sở dữ liệu
	config.ConnectDB()

	// Khởi tạo router
	r := gin.Default()

	// Cấu hình CORS tùy chỉnh
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // URL của frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true, // Cho phép gửi cookie hoặc thông tin xác thực
		MaxAge:           12 * time.Hour,
	}))

	// Thêm route để xử lý tất cả các yêu cầu OPTIONS
	r.OPTIONS("/*cors", func(c *gin.Context) {
		c.Status(204)
	})
	r.Static("/uploads", "./uploads") // Serve the uploads folder

	// Thiết lập các route
	routes.SetupRoutes(r)

	// Kết nối đến Redis
	err := config.ConnectRedis()
	if err != nil {
		log.Printf("Cảnh báo: Không thể kết nối đến Redis: %v", err)
	} else {
		defer config.RedisClient.Close()
		log.Println("Kết nối thành công đến Redis")
	}

	// Khởi động server trên cổng 8080
	log.Println("Khởi động server trên cổng :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Không thể khởi động server: %v", err)
	}
}
