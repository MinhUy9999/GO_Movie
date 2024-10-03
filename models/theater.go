package models

import "my-app/config"

// Theater struct represents the theater entity
type Theater struct {
	TheaterID int    `json:"theater_id"`
	Name      string `json:"name"`
	Location  string `json:"location"`
}

// CreateTheater inserts a new theater into the database
func CreateTheater(theater Theater) error {
	_, err := config.DB.Exec(
		"INSERT INTO THEATER (name, location) VALUES (?, ?)",
		theater.Name, theater.Location,
	)
	if err != nil {
		return err
	}
	return nil
}

// GetAllTheaters retrieves all theaters from the database
func GetAllTheaters() ([]Theater, error) {
	rows, err := config.DB.Query("SELECT theaterID, name, location FROM THEATER")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var theaters []Theater
	for rows.Next() {
		var theater Theater
		err := rows.Scan(&theater.TheaterID, &theater.Name, &theater.Location)
		if err != nil {
			return nil, err
		}
		theaters = append(theaters, theater)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return theaters, nil
}

// UpdateTheater updates the details of an existing theater in the database
func UpdateTheater(theater Theater) error {
	// Update the theater details based on the ID
	query := "UPDATE THEATER SET name = ?, location = ? WHERE theaterID = ?"
	_, err := config.DB.Exec(query, theater.Name, theater.Location, theater.TheaterID)
	return err
}

// DeleteTheater deletes a theater from the database
func DeleteTheater(theaterID int) error {
	_, err := config.DB.Exec("DELETE FROM THEATER WHERE theaterID = ?", theaterID)
	if err != nil {
		return err
	}
	return nil
}
