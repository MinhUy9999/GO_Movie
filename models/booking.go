package models

import (
	"database/sql"
	"errors"
	"my-app/config"
	"time"
)

// BookingDetails struct to represent booking data
type BookingDetails struct {
	BookingID   int64     `json:"booking_id"`
	UserID      int       `json:"user_id"`
	MovieID     int       `json:"movie_id"`
	ScreenID    int       `json:"screen_id"`
	BookingDate time.Time `json:"booking_date"`
	SeatsBooked int       `json:"seats_booked"`
	SeatIDs     []int     `json:"seat_ids"`
}

// BookSeats books seats for a user and a specific movie schedule
func BookSeats(userID, scheduleID int, seatIDs []int) (int64, error) {
	tx, err := config.DB.Begin()
	if err != nil {
		return 0, err
	}

	// Step 1: Check seat availability
	for _, seatID := range seatIDs {
		var isBooked bool
		err := tx.QueryRow("SELECT isBooked FROM SEAT WHERE seatID = ? AND screenID IN (SELECT screenID FROM SCHEDULE WHERE scheduleID = ?)", seatID, scheduleID).Scan(&isBooked)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		if isBooked {
			tx.Rollback()
			return 0, errors.New("one or more seats are already booked")
		}
	}

	// Step 2: Insert the booking record
	result, err := tx.Exec("INSERT INTO BOOKING (userID, movieID, screenID, bookingDate, seatsBooked) SELECT ?, movieID, screenID, ?, ? FROM SCHEDULE WHERE scheduleID = ?", userID, time.Now(), len(seatIDs), scheduleID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	bookingID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Step 3: Mark the seats as booked
	for _, seatID := range seatIDs {
		_, err := tx.Exec("UPDATE SEAT SET isBooked = true WHERE seatID = ?", seatID)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// Step 4: Commit the transaction
	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return bookingID, nil
}

// GetBookingDetailsByUserID retrieves all booking details by userID
// GetBookingDetailsByUserID retrieves all booking details by userID
func GetBookingDetailsByUserID(userID int) ([]BookingDetails, error) {
	var bookings []BookingDetails

	rows, err := config.DB.Query("SELECT bookingID, movieID, screenID, bookingDate, seatsBooked FROM BOOKING WHERE userID = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var booking BookingDetails
		var bookingDate sql.NullString // use NullString to handle NULLs gracefully

		booking.UserID = userID
		if err := rows.Scan(&booking.BookingID, &booking.MovieID, &booking.ScreenID, &bookingDate, &booking.SeatsBooked); err != nil {
			return nil, err
		}

		// Convert bookingDate to time.Time
		if bookingDate.Valid {
			booking.BookingDate, err = time.Parse("2006-01-02 15:04:05", bookingDate.String)
			if err != nil {
				return nil, err
			}
		}

		// Query booked seat IDs for each booking
		seatRows, err := config.DB.Query("SELECT seatID FROM SEAT WHERE isBooked = true AND screenID = ?", booking.ScreenID)
		if err != nil {
			return nil, err
		}
		defer seatRows.Close()

		for seatRows.Next() {
			var seatID int
			if err := seatRows.Scan(&seatID); err != nil {
				return nil, err
			}
			booking.SeatIDs = append(booking.SeatIDs, seatID)
		}

		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bookings, nil
}

// GetAllBookings retrieves all bookings
func GetAllBookings() ([]BookingDetails, error) {
	var bookings []BookingDetails

	rows, err := config.DB.Query("SELECT bookingID, userID, movieID, screenID, bookingDate, seatsBooked FROM BOOKING")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var booking BookingDetails
		var bookingDate sql.NullString // handle NULLs gracefully

		if err := rows.Scan(&booking.BookingID, &booking.UserID, &booking.MovieID, &booking.ScreenID, &bookingDate, &booking.SeatsBooked); err != nil {
			return nil, err
		}

		// Convert bookingDate to time.Time
		if bookingDate.Valid {
			booking.BookingDate, err = time.Parse("2006-01-02 15:04:05", bookingDate.String)
			if err != nil {
				return nil, err
			}
		}

		// Query booked seat IDs for each booking
		seatRows, err := config.DB.Query("SELECT seatID FROM SEAT WHERE isBooked = true AND screenID = ?", booking.ScreenID)
		if err != nil {
			return nil, err
		}
		defer seatRows.Close()

		for seatRows.Next() {
			var seatID int
			if err := seatRows.Scan(&seatID); err != nil {
				return nil, err
			}
			booking.SeatIDs = append(booking.SeatIDs, seatID)
		}

		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bookings, nil
}

// DeleteBooking deletes a booking by bookingID
func DeleteBooking(bookingID int64) error {
	result, err := config.DB.Exec("DELETE FROM BOOKING WHERE bookingID = ?", bookingID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no booking found with the given ID")
	}

	return nil
}
