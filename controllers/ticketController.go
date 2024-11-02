package controllers

import (
	"database/sql"
	"log"
	"my-app/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetTicketsByBookingIDHandler(c *gin.Context) {
	bookingIDStr := c.Param("bookingID")
	bookingID, err := strconv.Atoi(bookingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	tickets, err := models.GetTicketsByBookingID(bookingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tickets"})
		return
	}

	if len(tickets) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No tickets found for this booking ID"})
		return
	}

	// Convert sql.NullTime to string for JSON response
	type TicketResponse struct {
		TicketID  int     `json:"ticket_id"`
		BookingID int     `json:"booking_id"`
		SeatID    int     `json:"seat_id"`
		Fare      float64 `json:"fare"`
		IssuedAt  string  `json:"issued_at"`
		QRCode    string  `json:"qr_code"`
	}

	var responseTickets []TicketResponse
	for _, ticket := range tickets {
		issuedAt := ""
		if ticket.IssuedAt.Valid {
			issuedAt = ticket.IssuedAt.Time.Format(time.RFC3339)
		}
		responseTickets = append(responseTickets, TicketResponse{
			TicketID:  ticket.TicketID,
			BookingID: ticket.BookingID,
			SeatID:    ticket.SeatID,
			Fare:      ticket.Fare,
			IssuedAt:  issuedAt,
			QRCode:    ticket.QRCode,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": responseTickets})
}
func CreateTicketHandler(c *gin.Context) {
	var ticket models.Ticket

	if err := c.ShouldBindJSON(&ticket); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	exists, err := models.CheckBookingExists(ticket.BookingID)
	if err != nil {
		log.Printf("Error checking booking existence: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	// Set the IssuedAt field correctly
	ticket.IssuedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	ticketID, err := models.CreateTicket(&ticket)
	if err != nil {
		log.Printf("Error creating ticket: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ticket"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"ticket_id": ticketID})
}
