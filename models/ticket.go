// package models

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"time"

// 	"my-app/config"
// )

// type Ticket struct {
// 	TicketID  int          `json:"ticket_id"`
// 	BookingID int          `json:"booking_id"`
// 	SeatID    int          `json:"seat_id"`
// 	Fare      float64      `json:"fare"`
// 	IssuedAt  sql.NullTime `json:"issued_at"`
// 	QRCode    string       `json:"qr_code"`
// }

// // CheckBookingExists checks if a booking exists in the database
// func CheckBookingExists(bookingID int) (bool, error) {
// 	var exists bool
// 	query := "SELECT EXISTS(SELECT 1 FROM BOOKING WHERE bookingID = ?)"
// 	err := config.DB.QueryRow(query, bookingID).Scan(&exists)
// 	if err != nil {
// 		return false, fmt.Errorf("error checking booking existence: %v", err)
// 	}
// 	return exists, nil
// }

// // GetSeatsByBookingID retrieves seat IDs for a specific booking ID
// func GetSeatsByBookingID(bookingID int) ([]int, error) {
// 	var seatIDs []int
// 	query := `SELECT seatID FROM BOOKED_SEATS WHERE bookingID = ?`
// 	rows, err := config.DB.Query(query, bookingID)
// 	if err != nil {
// 		return nil, fmt.Errorf("error querying seats by booking ID: %v", err)
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var seatID int
// 		if err := rows.Scan(&seatID); err != nil {
// 			return nil, fmt.Errorf("error scanning seat ID: %v", err)
// 		}
// 		seatIDs = append(seatIDs, seatID)
// 	}

// 	if err = rows.Err(); err != nil {
// 		return nil, fmt.Errorf("error after scanning rows: %v", err)
// 	}

// 	return seatIDs, nil
// }

// // GetTicketsByBookingID retrieves tickets for a specific booking ID
// func GetTicketsByBookingID(bookingID int) ([]Ticket, error) {
// 	var tickets []Ticket
// 	query := `
//         SELECT ticketID, bookingID, seatID, fare, issuedAt, qrCode
//         FROM TICKET
//         WHERE bookingID = ?
//     `
// 	rows, err := config.DB.Query(query, bookingID)
// 	if err != nil {
// 		return nil, fmt.Errorf("error querying tickets: %v", err)
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var ticket Ticket
// 		err := rows.Scan(
// 			&ticket.TicketID,
// 			&ticket.BookingID,
// 			&ticket.SeatID,
// 			&ticket.Fare,
// 			&ticket.IssuedAt,
// 			&ticket.QRCode,
// 		)
// 		if err != nil {
// 			return nil, fmt.Errorf("error scanning ticket row: %v", err)
// 		}
// 		tickets = append(tickets, ticket)
// 	}

// 	if err = rows.Err(); err != nil {
// 		return nil, fmt.Errorf("error after scanning rows: %v", err)
// 	}

// 	return tickets, nil
// }

// // CreateTicket creates a ticket for a seat in a specific booking
// func CreateTicket(ticket *Ticket) (int, error) {
// 	query := `
//         INSERT INTO TICKET (seatID, bookingID, fare, issuedAt, qrCode)
//         VALUES (?, ?, ?, ?, ?)
//     `
// 	result, err := config.DB.Exec(query,
// 		ticket.SeatID,
// 		ticket.BookingID,
// 		ticket.Fare,
// 		sql.NullTime{Time: time.Now(), Valid: true},
// 		ticket.QRCode,
// 	)
// 	if err != nil {
// 		log.Printf("Error executing query in CreateTicket: %v", err)
// 		return 0, fmt.Errorf("error creating ticket: %v", err)
// 	}

// 	ticketID, err := result.LastInsertId()
// 	if err != nil {
// 		log.Printf("Error getting last insert ID in CreateTicket: %v", err)
// 		return 0, fmt.Errorf("error getting last insert ID: %v", err)
// 	}

//		return int(ticketID), nil
//	}
package models

import (
	"database/sql"
	"fmt"
	"time"

	"my-app/config"
)

type Booking struct {
	BookingID   int
	UserID      int
	MovieID     int
	ScreenID    int
	BookingDate time.Time
	SeatsBooked int
}

type Ticket struct {
	TicketID  int          `json:"ticket_id"`
	BookingID int          `json:"booking_id"`
	SeatID    int          `json:"seat_id"`
	Fare      float64      `json:"fare"`
	IssuedAt  sql.NullTime `json:"issued_at"`
	QRCode    string       `json:"qr_code"`
}

// GetBookingByID retrieves the booking details by bookingID
func GetBookingByID(bookingID int) (*Booking, error) {
	var booking Booking
	query := "SELECT bookingID, userID, movieID, screenID, bookingDate, seatsBooked FROM BOOKING WHERE bookingID = ?"
	err := config.DB.QueryRow(query, bookingID).Scan(
		&booking.BookingID,
		&booking.UserID,
		&booking.MovieID,
		&booking.ScreenID,
		&booking.BookingDate,
		&booking.SeatsBooked,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No booking found
		}
		return nil, fmt.Errorf("error retrieving booking: %v", err)
	}
	return &booking, nil
}

// CreateTicket inserts a new ticket record into the TICKET table
func CreateTicket(ticket *Ticket) (int, error) {
	query := `
        INSERT INTO TICKET (bookingID, seatID, fare, issuedAt, qrCode)
        VALUES (?, ?, ?, ?, ?)
    `
	result, err := config.DB.Exec(query,
		ticket.BookingID,
		ticket.SeatID,
		ticket.Fare,
		sql.NullTime{Time: time.Now(), Valid: true},
		ticket.QRCode,
	)
	if err != nil {
		return 0, fmt.Errorf("error creating ticket: %v", err)
	}

	ticketID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last insert ID: %v", err)
	}

	return int(ticketID), nil
}
