package routes

import (
	"my-app/controllers"
	"my-app/middlewares"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Ticket Booking API
// @version 1.0
// @description This is a sample server for movie ticket booking.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func SetupRoutes(r *gin.Engine) {
	// r := gin.Default()

	// Swagger documentation endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public routes
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	// User routes
	// User routes (chỉ cho vai trò "user")
	user := r.Group("/user")
	user.Use(middlewares.JWTAuthMiddleware("user"))
	{
		user.GET("/profile", controllers.UserProfile)
		user.PUT("/update/:id", controllers.UpdateUserByID)
		user.DELETE("/delete/:id", controllers.DeleteUserByID)

	}

	// Booking and Payment routes (cho phép cả "user" và "admin")
	booking := r.Group("/user")
	booking.Use(middlewares.JWTAuthMiddleware("user", "admin"))
	{
		booking.POST("/book", controllers.BookTickets)
		booking.POST("/payment", controllers.ProcessPayment)
	}
	// Ticket management routes
	ticket := r.Group("/tickets")
	{
		ticket.GET("/booking/:bookingID", controllers.GetTicketsByBookingIDHandler) // Lấy danh sách vé theo bookingID
		ticket.POST("", controllers.CreateTicketHandler)                            // Tạo vé mới
	}
	// Public routes for theaters
	theaterPublic := r.Group("/theater")
	{
		theaterPublic.GET("/", controllers.GetTheaters) // Get a list of all theaters
	}

	// Admin routes for theaters
	theaterAdmin := r.Group("/theater")
	theaterAdmin.Use(middlewares.JWTAuthMiddleware("admin")) // Apply JWT authentication for admin users
	{
		theaterAdmin.POST("/", controllers.CreateTheater)      // Create a new theater
		theaterAdmin.PUT("/:id", controllers.UpdateTheater)    // Update a theater by ID
		theaterAdmin.DELETE("/:id", controllers.DeleteTheater) // Delete a theater by ID
	}
	// Public routes for schedules
	schedulePublic := r.Group("/schedule")
	{
		schedulePublic.GET("/screen/:screenID", controllers.GetSchedulesByScreenIDHandler) // Get schedules by screenID
		schedulePublic.GET("/:id", controllers.GetScheduleByIDHandler)
	}

	// Admin routes for schedules
	scheduleAdmin := r.Group("/schedule")
	scheduleAdmin.Use(middlewares.JWTAuthMiddleware("admin"))
	{
		scheduleAdmin.POST("/", controllers.CreateScheduleHandler) // Create a new schedule
		scheduleAdmin.GET("/", controllers.GetSchedulesHandler)    // Get all schedules
		// Get schedule information by ID
		scheduleAdmin.PUT("/:id", controllers.UpdateScheduleHandler)    // Update schedule by ID
		scheduleAdmin.DELETE("/:id", controllers.DeleteScheduleHandler) // Delete schedule by ID
	}

	// Seat management routes (for all authenticated users)
	seats := r.Group("/seats")
	seats.Use(middlewares.JWTAuthMiddleware()) // Dành cho người dùng đã đăng nhập
	{
		seats.GET("/:screenID", controllers.GetSeatsByScreenID)
	}

	// Admin-only routes for seats
	seatsAdmin := r.Group("/seats")
	seatsAdmin.Use(middlewares.JWTAuthMiddleware("admin"))
	{
		seatsAdmin.POST("/", controllers.CreateSeat)
	}

	// Room management routes (Admin only)
	room := r.Group("/room")
	{
		// Public endpoint for getting rooms by theater ID
		room.GET("/:theaterID", controllers.GetRoomsByTheaterID)
	}

	// Admin-only routes for rooms
	roomAdmin := r.Group("/room")
	roomAdmin.Use(middlewares.JWTAuthMiddleware("admin"))
	{
		roomAdmin.POST("/", controllers.CreateRoom)
		roomAdmin.PUT("/:id", controllers.UpdateRoom)
		roomAdmin.DELETE("/:id", controllers.DeleteRoom)
	}
	// Screen management routes (Admin only)
	// Screen management routes
	screen := r.Group("/screen")
	{
		// Public endpoint for getting screens by room ID
		screen.GET("/:roomID", controllers.GetScreensByRoomID)
	}

	// Admin-only routes for screens
	screenAdmin := r.Group("/screen")
	screenAdmin.Use(middlewares.JWTAuthMiddleware("admin"))
	{
		screenAdmin.POST("/", controllers.CreateScreen)
		screenAdmin.PUT("/:id", controllers.UpdateScreen)
		screenAdmin.DELETE("/:id", controllers.DeleteScreen)
	}
	// Chat routes
	r.GET("/chat/messages", controllers.GetMessagesHandler) // Lấy tất cả tin nhắn
	r.POST("/chat/messages", controllers.AddMessageHandler) // Thêm tin nhắn mới
	r.GET("/ws", controllers.ChatWebSocketHandler)
	// Movie management routes (Admin only)
	moviePublic := r.Group("/movie")
	{
		moviePublic.GET("/", controllers.GetAllMovies)    // Lấy danh sách tất cả phim
		moviePublic.GET("/:id", controllers.GetMovieByID) // Lấy thông tin chi tiết phim
	}

	// Định nghĩa route cho các endpoint yêu cầu quyền admin
	movieAdmin := r.Group("/movie")
	movieAdmin.Use(middlewares.JWTAuthMiddleware("admin"))
	{
		movieAdmin.POST("/", controllers.CreateMovie)
		movieAdmin.PUT("/:id", controllers.UpdateMovie)
		movieAdmin.DELETE("/:id", controllers.DeleteMovie)
	}

	// Admin-only routes
	admin := r.Group("/admin")
	admin.Use(middlewares.JWTAuthMiddleware("admin"))
	{
		admin.GET("/users", controllers.GetAllUsers)
	}

}
