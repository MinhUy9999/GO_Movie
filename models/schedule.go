package models

import (
	"fmt"
	"my-app/config"
)

// Schedule - Model for movie schedule
type Schedule struct {
	ScheduleID     int     `json:"scheduleID"`     // Unique identifier for the schedule
	MovieID        int     `json:"movieID"`        // Unique identifier for the movie
	ScreenID       int     `json:"screenID"`       // Unique identifier for the screen
	ShowTime       string  `json:"showTime"`       // Show time of the movie
	AvailableSeats int     `json:"availableSeats"` // Number of available seats
	Fare           float64 `json:"fare"`           // Price per seat for the schedule
}

// CreateSchedule - Add a new schedule
func CreateSchedule(schedule *Schedule) error {
	query := "INSERT INTO SCHEDULE (movieID, screenID, showTime, availableSeats, fare) VALUES (?, ?, ?, ?, ?)"
	result, err := config.DB.Exec(query, schedule.MovieID, schedule.ScreenID, schedule.ShowTime, schedule.AvailableSeats, schedule.Fare)
	if err != nil {
		return err
	}

	lastInsertID, err := result.LastInsertId()
	if err == nil {
		schedule.ScheduleID = int(lastInsertID)
	}
	return err
}

// GetSchedules - Retrieve the list of schedules
func GetSchedules() ([]Schedule, error) {
	query := "SELECT scheduleID, movieID, screenID, showTime, availableSeats, fare FROM SCHEDULE"
	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []Schedule
	for rows.Next() {
		var schedule Schedule
		if err := rows.Scan(&schedule.ScheduleID, &schedule.MovieID, &schedule.ScreenID, &schedule.ShowTime, &schedule.AvailableSeats, &schedule.Fare); err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}

	return schedules, nil
}

// GetScheduleByID - Retrieve a schedule by ID
func GetScheduleByID(id int) (Schedule, error) {
	query := "SELECT scheduleID, movieID, screenID, showTime, availableSeats, fare FROM SCHEDULE WHERE scheduleID = ?"
	var schedule Schedule
	err := config.DB.QueryRow(query, id).Scan(&schedule.ScheduleID, &schedule.MovieID, &schedule.ScreenID, &schedule.ShowTime, &schedule.AvailableSeats, &schedule.Fare)
	return schedule, err
}

// UpdateSchedule - Update a schedule
func UpdateSchedule(schedule *Schedule) error {
	query := `
        UPDATE SCHEDULE
        SET movieID = ?, screenID = ?, showTime = ?, availableSeats = ?, fare = ?
        WHERE scheduleID = ?
    `

	// Log the values to be used in the query
	fmt.Printf("Updating schedule with values: MovieID=%d, ScreenID=%d, ShowTime=%s, AvailableSeats=%d, Fare=%.2f, ScheduleID=%d\n",
		schedule.MovieID, schedule.ScreenID, schedule.ShowTime, schedule.AvailableSeats, schedule.Fare, schedule.ScheduleID)

	_, err := config.DB.Exec(query, schedule.MovieID, schedule.ScreenID, schedule.ShowTime, schedule.AvailableSeats, schedule.Fare, schedule.ScheduleID)
	return err
}

// DeleteSchedule - Delete a schedule
func DeleteSchedule(id int) error {
	query := "DELETE FROM SCHEDULE WHERE scheduleID = ?"
	_, err := config.DB.Exec(query, id)
	return err
}

// GetSchedulesByScreenID - Retrieve schedules by screenID
func GetSchedulesByScreenID(screenID int) ([]Schedule, error) {
	query := "SELECT scheduleID, movieID, screenID, showTime, availableSeats, fare FROM SCHEDULE WHERE screenID = ?"
	rows, err := config.DB.Query(query, screenID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []Schedule
	for rows.Next() {
		var schedule Schedule
		if err := rows.Scan(&schedule.ScheduleID, &schedule.MovieID, &schedule.ScreenID, &schedule.ShowTime, &schedule.AvailableSeats, &schedule.Fare); err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}

	return schedules, nil
}
