package models // Ensure the package is declared
import "my-app/config"

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
