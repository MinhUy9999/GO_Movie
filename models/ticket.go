package models

import (
	"database/sql"
	"fmt"
	"time"

	"my-app/config"
)

type Ticket struct {
	TicketID  int          `json:"ticket_id"`
	BookingID int          `json:"booking_id"`
	SeatID    int          `json:"seat_id"`
	Fare      float64      `json:"fare"`
	IssuedAt  sql.NullTime `json:"issued_at"`
	QRCode    string       `json:"qr_code"`
}

func CheckBookingExists(bookingID int) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM BOOKING WHERE bookingID = ?)"
	err := config.DB.QueryRow(query, bookingID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking booking existence: %v", err)
	}
	return exists, nil
}

func GetTicketsByBookingID(bookingID int) ([]Ticket, error) {
	var tickets []Ticket
	query := `
        SELECT ticketID, bookingID, seatID, fare, issuedAt, qrCode 
        FROM TICKET 
        WHERE bookingID = ?
    `
	rows, err := config.DB.Query(query, bookingID)
	if err != nil {
		return nil, fmt.Errorf("error querying tickets: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var ticket Ticket
		err := rows.Scan(
			&ticket.TicketID,
			&ticket.BookingID,
			&ticket.SeatID,
			&ticket.Fare,
			&ticket.IssuedAt,
			&ticket.QRCode,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning ticket row: %v", err)
		}
		tickets = append(tickets, ticket)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning rows: %v", err)
	}

	return tickets, nil
}

func CreateTicket(ticket *Ticket) (int, error) {
	query := `
        INSERT INTO TICKET (seatID, bookingID, fare, issuedAt, qrCode)
        VALUES (?, ?, ?, ?, ?)
    `
	result, err := config.DB.Exec(query,
		ticket.SeatID,
		ticket.BookingID,
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
