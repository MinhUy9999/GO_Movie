package models

import (
	"errors"
	"my-app/config"
)

type Payment struct {
	BookingID     int     `json:"booking_id"`
	Amount        float64 `json:"amount"`
	PaymentStatus string  `json:"payment_status"` // PAID or PENDING
}

// ProcessPayment processes the payment for a booking
func ProcessPayment(userID int, payment Payment) error {
	// Check if the booking exists and belongs to the user
	var bookingExists bool
	err := config.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM BOOKING WHERE bookingID = ? AND userID = ?)", payment.BookingID, userID,
	).Scan(&bookingExists)

	if err != nil {
		return err
	}
	if !bookingExists {
		return errors.New("invalid booking ID or booking does not belong to the user")
	}

	// Update the payment status in the database
	_, err = config.DB.Exec(
		"INSERT INTO PAYMENT (bookingID, amount, paymentStatus) VALUES (?, ?, ?)",
		payment.BookingID, payment.Amount, payment.PaymentStatus,
	)
	if err != nil {
		return err
	}

	// Optionally, update the booking status to reflect that payment has been made
	if payment.PaymentStatus == "PAID" {
		_, err = config.DB.Exec("UPDATE BOOKING SET paymentStatus = ? WHERE bookingID = ?", "PAID", payment.BookingID)
		if err != nil {
			return err
		}
	}

	return nil
}
