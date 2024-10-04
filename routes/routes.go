package routes

import (
	"my-app/controllers"
	"my-app/middlewares"

	"github.com/gin-gonic/gin"
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
func SetupRoutes() *gin.Engine {
	r := gin.Default()

	// Use Redis middleware
	r.Use(middlewares.UseRedis())

	// Swagger documentation endpoint

	// Public routes
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	// User routes
	user := r.Group("/user")
	user.Use(middlewares.JWTAuthMiddleware("user"))
	{
		user.GET("/profile", controllers.UserProfile)
		user.PUT("/update/:id", controllers.UpdateUserByID)
		user.DELETE("/delete/:id", controllers.DeleteUserByID)
		user.POST("/payment", controllers.ProcessPayment)
		user.POST("/book", controllers.BookTickets)
	}

	// Theater management routes (Admin only)
	theater := r.Group("/theater")
	theater.Use(middlewares.JWTAuthMiddleware("admin"))
	{
		theater.POST("/", controllers.CreateTheater)
		theater.GET("/", controllers.GetTheaters)
		theater.PUT("/:id", controllers.UpdateTheater)
		theater.DELETE("/:id", controllers.DeleteTheater)
	}

	// Seat management routes (for all authenticated users)
	seats := r.Group("/seats")
	seats.Use(middlewares.JWTAuthMiddleware(""))
	{
		seats.GET("/:screenID", controllers.GetSeatsByScreenID)
	}

	// Room management routes (Admin only)
	room := r.Group("/room")
	room.Use(middlewares.JWTAuthMiddleware("admin"))
	{
		room.POST("/", controllers.CreateRoom)
		room.GET("/:theaterID", controllers.GetRoomsByTheaterID)
		room.PUT("/:id", controllers.UpdateRoom)
		room.DELETE("/:id", controllers.DeleteRoom)
	}

	// Movie management routes (Admin only)
	movie := r.Group("/movie")
	movie.Use(middlewares.JWTAuthMiddleware("admin"))
	{
		movie.POST("/", controllers.CreateMovie)
		movie.GET("/", controllers.GetAllMovies)
		movie.PUT("/:id", controllers.UpdateMovie)
		movie.DELETE("/:id", controllers.DeleteMovie)
	}

	// Admin-only routes
	admin := r.Group("/admin")
	admin.Use(middlewares.JWTAuthMiddleware("admin"))
	{
		admin.GET("/users", controllers.GetAllUsers)
	}

	return r
}
