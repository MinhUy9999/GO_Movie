package controllers

import (
	"log"
	"my-app/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// ChatResponse đại diện cho cấu trúc phản hồi cho các tin nhắn chat
type ChatResponse struct {
	Success bool                 `json:"success"`
	Message string               `json:"message,omitempty"`
	Data    []models.ChatMessage `json:"data,omitempty"`
}

// Biến upgrader dùng để nâng cấp kết nối HTTP thành WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Cho phép mọi nguồn gốc kết nối, điều chỉnh nếu cần hạn chế nguồn gốc
	},
}

// GetMessagesHandler xử lý các yêu cầu GET để truy xuất tất cả tin nhắn chat
func GetMessagesHandler(c *gin.Context) {
	// Đặt header cho phản hồi
	c.Header("Content-Type", "application/json")

	// Lấy tin nhắn từ model
	messages, err := models.GetMessages()
	if err != nil {
		log.Printf("Error retrieving messages: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
		return
	}

	response := ChatResponse{
		Success: true,
		Data:    messages,
	}

	// Trả về phản hồi dạng JSON
	c.JSON(http.StatusOK, response)
}

// AddMessageHandler xử lý các yêu cầu POST để thêm tin nhắn mới vào cơ sở dữ liệu
func AddMessageHandler(c *gin.Context) {
	var reqBody struct {
		UserID      int    `json:"userID,omitempty"`
		MessageText string `json:"messageText" binding:"required"`
	}

	// Phân tích yêu cầu JSON
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Lưu tin nhắn vào cơ sở dữ liệu
	err := models.AddMessage(reqBody.UserID, reqBody.MessageText)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save message"})
		return
	}

	response := ChatResponse{
		Success: true,
		Message: "Message saved successfully",
	}

	// Trả về phản hồi dạng JSON
	c.JSON(http.StatusCreated, response)
}

// ChatWebSocketHandler xử lý các kết nối WebSocket cho chat thời gian thực
func ChatWebSocketHandler(c *gin.Context) {
	// Nâng cấp kết nối HTTP thành WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade to WebSocket"})
		return
	}
	defer conn.Close()

	// Lắng nghe tin nhắn từ client
	for {
		var msg models.ChatMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		// Lưu tin nhắn vào cơ sở dữ liệu
		err = models.AddMessage(msg.UserID, msg.MessageText)
		if err != nil {
			log.Println("Failed to save message:", err)
			break
		}

		// Phát lại tin nhắn cho client (Echo message)
		if err := conn.WriteJSON(msg); err != nil {
			log.Println("Error writing message:", err)
			break
		}
	}
}
