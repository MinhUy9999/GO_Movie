package models

import "my-app/config"

type Room struct {
	RoomID     int `json:"room_id"`
	TheaterID  int `json:"theater_id"`
	RoomNumber int `json:"room_number"`
}

func CreateRoom(room Room) error {
	_, err := config.DB.Exec(
		"INSERT INTO ROOM (theaterID, roomNumber) VALUES (?, ?)",
		room.TheaterID, room.RoomNumber,
	)
	if err != nil {
		return err
	}
	return nil
}
func GetRoomsByTheaterID(theaterID int) ([]Room, error) {
	rows, err := config.DB.Query("SELECT roomID, theaterID, roomNumber FROM ROOM WHERE theaterID = ?", theaterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []Room
	for rows.Next() {
		var room Room
		err := rows.Scan(&room.RoomID, &room.TheaterID, &room.RoomNumber)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return rooms, nil
}
func UpdateRoom(room Room) error {
	_, err := config.DB.Exec(
		"UPDATE ROOM SET theaterID = ?, roomNumber = ? WHERE roomID = ?",
		room.TheaterID, room.RoomNumber, room.RoomID,
	)
	if err != nil {
		return err
	}
	return nil
}
func DeleteRoom(roomID int) error {
	_, err := config.DB.Exec("DELETE FROM ROOM WHERE roomID = ?", roomID)
	if err != nil {
		return err
	}
	return nil
}
