package models // Ensure the package is declared

import "my-app/config"


type Screen struct {
    ScreenID     int     `json:"screen_id"`
    RoomID       int     `json:"room_id"`
    ScreenNumber int     `json:"screen_number"`
    Opacity      float64 `json:"opacity"`
}

// CreateScreen adds a new screen to the database
func CreateScreen(screen *Screen) error {
    query := "INSERT INTO SCREEN (roomID, screenNumber, opacity) VALUES (?, ?, ?)"
    result, err := config.DB.Exec(query, screen.RoomID, screen.ScreenNumber, screen.Opacity)
    if err != nil {
        return err
    }

    screenID, err := result.LastInsertId()
    if err != nil {
        return err
    }

    screen.ScreenID = int(screenID)
    return nil
}

// GetScreensByRoomID retrieves all screens for a specific room ID
func GetScreensByRoomID(roomID int) ([]Screen, error) {
    rows, err := config.DB.Query("SELECT screenID, roomID, screenNumber, opacity FROM SCREEN WHERE roomID = ?", roomID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var screens []Screen
    for rows.Next() {
        var screen Screen
        err := rows.Scan(&screen.ScreenID, &screen.RoomID, &screen.ScreenNumber, &screen.Opacity)
        if err != nil {
            return nil, err
        }
        screens = append(screens, screen)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return screens, nil
}

// UpdateScreen modifies an existing screen in the database
func UpdateScreen(screen *Screen) error {
    query := "UPDATE SCREEN SET roomID = ?, screenNumber = ?, opacity = ? WHERE screenID = ?"
    _, err := config.DB.Exec(query, screen.RoomID, screen.ScreenNumber, screen.Opacity, screen.ScreenID)
    return err
}

// DeleteScreen removes a screen from the database based on its ID
func DeleteScreen(screenID int) error {
    query := "DELETE FROM SCREEN WHERE screenID = ?"
    _, err := config.DB.Exec(query, screenID)
    return err
}
