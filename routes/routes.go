package routes

import (
	"my-app/controllers"
	_ "my-app/docs"
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
func SetupRoutes() *gin.Engine {
	r := gin.Default()

	// Swagger documentation endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Auth routes
	// @Summary Register a new user
	// @Tags Authentication
	// @Accept  json
	// @Produce  json
	// @Param   user body models.User true "User registration details"
	// @Success 200 {object} map[string]interface{}
	// @Failure 400 {object} map[string]interface{}
	// @Router /register [post]
	r.POST("/register", controllers.Register)

	// @Summary Login a user
	// @Tags Authentication
	// @Accept  json
	// @Produce  json
	// @Param   login body models.User true "User login details"
	// @Success 200 {object} map[string]interface{}
	// @Failure 400 {object} map[string]interface{}
	// @Router /login [post]
	r.POST("/login", controllers.Login)

	// Protected routes for authenticated users
	auth := r.Group("/user")
	auth.Use(middlewares.JWTAuthMiddleware("")) // Empty string means any authenticated user
	{
		// @Summary Get user profile
		// @Tags Users
		// @Produce  json
		// @Success 200 {object} map[string]interface{}
		// @Router /user/profile [get]
		auth.GET("/profile", controllers.UserProfile)

		// @Summary Update user details
		// @Tags Users
		// @Accept  json
		// @Produce  json
		// @Param   user body models.User true "Updated user details"
		// @Success 200 {object} map[string]interface{}
		// @Failure 400 {object} map[string]interface{}
		// @Router /user/update [put]
		auth.PUT("/update/:id", controllers.UpdateUserByID)

		// @Summary Delete user account
		// @Tags Users
		// @Produce  json
		// @Success 200 {object} map[string]interface{}
		// @Router /user/delete [delete]
		auth.DELETE("/delete/:id", controllers.DeleteUserByID)

		// @Summary Process a payment
		// @Tags Payments
		// @Accept  json
		// @Produce  json
		// @Param   payment body models.Payment true "Payment details"
		// @Success 200 {object} map[string]interface{}
		// @Failure 400 {object} map[string]interface{}
		// @Router /user/payment [post]
		auth.POST("/payment", controllers.ProcessPayment)

		// @Summary Book tickets
		// @Tags Booking
		// @Accept  json
		// @Produce  json
		// @Param   booking body models.Booking true "Booking details"
		// @Success 200 {object} map[string]interface{}
		// @Failure 400 {object} map[string]interface{}
		// @Router /user/book [post]
		auth.POST("/book", controllers.BookTickets)
	}

	// Theater management routes (Admin only)
	theater := r.Group("/theater")
	theater.Use(middlewares.JWTAuthMiddleware("admin")) // Only admins can manage theaters
	{
		// @Summary Create a new theater
		// @Tags Theater
		// @Accept  json
		// @Produce  json
		// @Param   theater body models.Theater true "Theater details"
		// @Success 200 {object} map[string]interface{}
		// @Failure 400 {object} map[string]interface{}
		// @Router /theater [post]
		theater.POST("/", controllers.CreateTheater)

		// @Summary Get list of all theaters
		// @Tags Theater
		// @Produce  json
		// @Success 200 {object} map[string]interface{}
		// @Failure 500 {object} map[string]interface{}
		// @Router /theater [get]
		theater.GET("/", controllers.GetTheaters)

		// @Summary Update a theater by ID
		// @Tags Theater
		// @Accept  json
		// @Produce  json
		// @Param   id path int true "Theater ID"
		// @Param   theater body models.Theater true "Updated theater details"
		// @Success 200 {object} map[string]interface{}
		// @Failure 400 {object} map[string]interface{}
		// @Router /theater/{id} [put]
		theater.PUT("/:id", controllers.UpdateTheater)

		// @Summary Delete a theater by ID
		// @Tags Theater
		// @Produce  json
		// @Param   id path int true "Theater ID"
		// @Success 200 {object} map[string]interface{}
		// @Failure 400 {object} map[string]interface{}
		// @Router /theater/{id} [delete]
		theater.DELETE("/:id", controllers.DeleteTheater)
	}

	// Seat management routes (for all users)
	// @Summary Fetch seats by screen ID
	// @Tags Seat
	// @Produce  json
	// @Param   screenID path int true "Screen ID"
	// @Success 200 {object} map[string]interface{}
	// @Failure 500 {object} map[string]interface{}
	// @Router /seats/{screenID} [get]
	r.GET("/seats/:screenID", controllers.GetSeatsByScreenID)

	// Room management routes (Admin only)
	room := r.Group("/room")
	room.Use(middlewares.JWTAuthMiddleware("admin")) // Only admins can manage rooms
	{
		// @Summary Create a new room
		// @Tags Room
		// @Accept  json
		// @Produce  json
		// @Param   room body models.Room true "Room details"
		// @Success 200 {object} map[string]interface{}
		// @Failure 400 {object} map[string]interface{}
		// @Router /room [post]
		room.POST("/", controllers.CreateRoom)

		// @Summary Get rooms by theater ID
		// @Tags Room
		// @Produce  json
		// @Param   theaterID path int true "Theater ID"
		// @Success 200 {object} map[string]interface{}
		// @Failure 500 {object} map[string]interface{}
		// @Router /room/{theaterID} [get]
		room.GET("/:theaterID", controllers.GetRoomsByTheaterID)

		// @Summary Update a room by ID
		// @Tags Room
		// @Accept  json
		// @Produce  json
		// @Param   id path int true "Room ID"
		// @Param   room body models.Room true "Updated room details"
		// @Success 200 {object} map[string]interface{}
		// @Failure 400 {object} map[string]interface{}
		// @Router /room/{id} [put]
		room.PUT("/:id", controllers.UpdateRoom)

		// @Summary Delete a room by ID
		// @Tags Room
		// @Produce  json
		// @Param   id path int true "Room ID"
		// @Success 200 {object} map[string]interface{}
		// @Failure 400 {object} map[string]interface{}
		// @Router /room/{id} [delete]
		room.DELETE("/:id", controllers.DeleteRoom)
	}

	// Movie management routes (Admin only)
	movie := r.Group("/movie")
	movie.Use(middlewares.JWTAuthMiddleware("admin")) // Only admins can manage movies
	{
		// @Summary Create a new movie
		// @Tags Movie
		// @Accept  json
		// @Produce  json
		// @Param   movie body models.Movie true "Movie details"
		// @Success 200 {object} map[string]interface{}
		// @Failure 400 {object} map[string]interface{}
		// @Router /movie [post]
		movie.POST("/", controllers.CreateMovie)

		// @Summary Get list of all movies
		// @Tags Movie
		// @Produce  json
		// @Success 200 {object} map[string]interface{}
		// @Failure 500 {object} map[string]interface{}
		// @Router /movie [get]
		movie.GET("/", controllers.GetAllMovies)

		// @Summary Update a movie by ID
		// @Tags Movie
		// @Accept  json
		// @Produce  json
		// @Param   id path int true "Movie ID"
		// @Param   movie body models.Movie true "Updated movie details"
		// @Success 200 {object} map[string]interface{}
		// @Failure 400 {object} map[string]interface{}
		// @Router /movie/{id} [put]
		movie.PUT("/:id", controllers.UpdateMovie)

		// @Summary Delete a movie by ID
		// @Tags Movie
		// @Produce  json
		// @Param   id path int true "Movie ID"
		// @Success 200 {object} map[string]interface{}
		// @Failure 400 {object} map[string]interface{}
		// @Router /movie/{id} [delete]
		movie.DELETE("/:id", controllers.DeleteMovie)
	}

	// Admin-only routes
	admin := r.Group("/admin")
	admin.Use(middlewares.JWTAuthMiddleware("admin")) // Only admins can access these routes
	{
		// @Summary Get all users (Admin only)
		// @Tags Admin
		// @Produce  json
		// @Success 200 {object} map[string]interface{}
		// @Failure 500 {object} map[string]interface{}
		// @Router /admin/users [get]
		admin.GET("/users", controllers.GetAllUsers)
	}

	return r
}
