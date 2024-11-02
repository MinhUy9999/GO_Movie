package controllers

import (
	"fmt"
	"io"
	"math/rand"
	"my-app/models"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Create a new movie (C)
func CreateMovie(c *gin.Context) {
	var movie models.Movie

	// Bind các trường text (title, genre, duration) - không bind file
	if err := c.ShouldBind(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": err.Error(),
		})
		return
	}

	// Xử lý upload file hình ảnh thủ công
	file, fileHeader, err := c.Request.FormFile("picture")
	if err == nil {
		defer file.Close()

		// Tạo tên file độc nhất
		rand.Seed(time.Now().UnixNano())
		ext := filepath.Ext(fileHeader.Filename)
		if ext == "" {
			ext = ".png" // Đặt mặc định nếu không có đuôi mở rộng
		}
		filename := fmt.Sprintf("movie_%d_%d%s", time.Now().Unix(), rand.Int(), ext)
		filePath := filepath.Join("uploads", "images", filename)
		movie.Picture = filepath.Join("images", filename)
		// Đảm bảo thư mục tồn tại
		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
			return
		}

		// Tạo file mới trên server
		dst, err := os.Create(filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
			return
		}
		defer dst.Close()

		// Copy nội dung từ file upload vào file mới
		if _, err = io.Copy(dst, file); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}

		// Lưu đường dẫn file ảnh vào trường Picture của movie
		movie.Picture = filePath
	} else {
		movie.Picture = "" // Nếu không có file, đặt rỗng
	}

	// Lưu movie vào database
	err = models.CreateMovie(movie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Movie created successfully"})
}

// Get movie details or list of movies (R)
func GetAllMovies(c *gin.Context) {
	movies, err := models.GetAllMovies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"movies": movies})
}
func GetMovieByID(c *gin.Context) {
	// Lấy movieID từ URL parameter
	movieIDStr := c.Param("id")
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	// Gọi hàm GetMovieByID trong models để lấy thông tin phim
	movie, err := models.GetMovieByID(movieID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	// Trả về thông tin phim dưới dạng JSON
	c.JSON(http.StatusOK, gin.H{"movie": movie})
}

// Update movie details (U)
func UpdateMovie(c *gin.Context) {
	var movie models.Movie

	// Get the movie ID from the URL parameter
	movieIDStr := c.Param("id")
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}
	movie.MovieID = movieID

	// Bind form data
	if err := c.ShouldBind(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Retrieve the existing movie data
	existingMovie, err := models.GetMovieByID(movie.MovieID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Movie not found"})
		return
	}

	// Handle image file upload
	file, fileHeader, err := c.Request.FormFile("picture")
	if err == nil {
		defer file.Close()

		// Generate a unique filename
		rand.Seed(time.Now().UnixNano())
		ext := path.Ext(fileHeader.Filename)
		if ext == "" {
			ext = ".png" // Default extension
		}
		filename := fmt.Sprintf("movie_%d_%d%s", time.Now().Unix(), rand.Int(), ext)
		filePath := filepath.Join("uploads", "images", filename)

		// Ensure the directory exists
		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
			return
		}

		// Create a new file on the server
		dst, err := os.Create(filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
			return
		}
		defer dst.Close()

		// Copy the contents from the uploaded file to the new file
		if _, err = io.Copy(dst, file); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}

		// Delete the old image if it exists
		if existingMovie.Picture != "" {
			os.Remove(existingMovie.Picture)
		}

		// Update the Picture field to store the new file path
		movie.Picture = filePath
	} else if err != http.ErrMissingFile {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file upload"})
		return
	} else {
		// Keep the existing picture if no new image is provided
		movie.Picture = existingMovie.Picture
	}

	err = models.UpdateMovie(movie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Movie updated successfully"})
}

// Delete movie (D)
func DeleteMovie(c *gin.Context) {
	movieIDStr := c.Param("id") // Get the movie ID from the request parameters (string)

	// Convert the movieID string to an integer
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	// Retrieve the existing movie data to get the image path
	existingMovie, err := models.GetMovieByID(movieID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Movie not found"})
		return
	}

	// Call the models.DeleteMovie function with the converted integer movieID
	err = models.DeleteMovie(movieID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Delete the associated image file if it exists
	if existingMovie.Picture != "" {
		os.Remove(existingMovie.Picture)
	}

	c.JSON(http.StatusOK, gin.H{"status": "Movie deleted successfully"})
}
