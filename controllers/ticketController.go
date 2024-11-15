package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"my-app/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateTicketsForBookingHandler(c *gin.Context) {
	var req struct {
		BookingID int `json:"booking_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	// Retrieve booking details
	booking, err := models.GetBookingByID(req.BookingID)
	if err != nil {
		log.Printf("Database error while retrieving booking: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if booking == nil {
		log.Printf("Invalid booking ID: %d", req.BookingID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	// Create tickets based on seatsBooked
	var tickets []models.Ticket
	for i := 1; i <= booking.SeatsBooked; i++ {
		ticket := models.Ticket{
			BookingID: booking.BookingID,
			SeatID:    i,      // Increment seat ID for each ticket
			Fare:      1000.0, // Set fare as needed or retrieve dynamically
			IssuedAt:  sql.NullTime{Time: time.Now(), Valid: true},
			QRCode:    fmt.Sprintf("QR_%d_%d", booking.BookingID, i),
		}

		// Insert ticket into database
		ticketID, err := models.CreateTicket(&ticket)
		if err != nil {
			log.Printf("Error creating ticket for seat ID %d: %v", i, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tickets"})
			return
		}
		ticket.TicketID = ticketID
		tickets = append(tickets, ticket)
	}

	// Return the created tickets
	log.Printf("Successfully created %d tickets for booking ID %d", booking.SeatsBooked, booking.BookingID)
	c.JSON(http.StatusCreated, gin.H{"tickets": tickets})
}
func GetTicketsByBookingIDHandler(c *gin.Context) {
	bookingID, err := strconv.Atoi(c.Param("bookingID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	tickets, err := models.GetBookingByID(bookingID)
	if err != nil {
		log.Printf("Error retrieving tickets for booking ID %d: %v", bookingID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tickets"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tickets": tickets})
}
