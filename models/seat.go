package models // Ensure the package is declared
import (
	"fmt"
	"my-app/config"
)

type Seat struct {
	SeatID     int  `json:"seat_id"`
	ScreenID   int  `json:"screen_id"`
	SeatNumber int  `json:"seat_number"`
	IsBooked   bool `json:"is_booked"`
}

// CreateSeat adds a new seat to the database
func CreateSeat(seat *Seat) error {
	query := "INSERT INTO SEAT (screenID, seatNumber, isBooked) VALUES (?, ?, ?)"
	result, err := config.DB.Exec(query, seat.ScreenID, seat.SeatNumber, seat.IsBooked)
	if err != nil {
		return err
	}

	seatID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	seat.SeatID = int(seatID)
	return nil
}

// GetAllSeats retrieves all seats from the database
func GetAllSeats() ([]Seat, error) {
	rows, err := config.DB.Query("SELECT seatID, screenID, seatNumber, isBooked FROM SEAT")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seats []Seat
	for rows.Next() {
		var seat Seat
		err := rows.Scan(&seat.SeatID, &seat.ScreenID, &seat.SeatNumber, &seat.IsBooked)
		if err != nil {
			return nil, err
		}
		seats = append(seats, seat)
	}

	return seats, nil
}

// UpdateSeat updates an existing seat
// UpdateSeat updates an existing seat in the database
func UpdateSeat(seat *Seat) error {
	query := "UPDATE SEAT SET screenID = ?, seatNumber = ?, isBooked = ? WHERE seatID = ?"
	result, err := config.DB.Exec(query, seat.ScreenID, seat.SeatNumber, seat.IsBooked, seat.SeatID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows were updated; seatID %d might not exist", seat.SeatID)
	}

	return nil
}

// DeleteSeat removes a seat from the database by its ID
func DeleteSeat(seatID int) error {
	query := "DELETE FROM SEAT WHERE seatID = ?"
	_, err := config.DB.Exec(query, seatID)
	return err
}

// Example function to get seats by screen ID
func GetSeatsByScreenID(screenID int) ([]Seat, error) {
	rows, err := config.DB.Query("SELECT seatID, screenID, seatNumber, isBooked FROM SEAT WHERE screenID = ?", screenID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seats []Seat
	for rows.Next() {
		var seat Seat
		err := rows.Scan(&seat.SeatID, &seat.ScreenID, &seat.SeatNumber, &seat.IsBooked)
		if err != nil {
			return nil, err
		}
		seats = append(seats, seat)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return seats, nil
}
