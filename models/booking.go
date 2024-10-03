package models

import (
	"errors"
	"my-app/config"
	"time"
)

// BookSeats books seats for a user and a specific movie schedule
func BookSeats(userID, scheduleID int, seatIDs []int) (int64, error) {
	// Start a database transaction to ensure atomicity
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
