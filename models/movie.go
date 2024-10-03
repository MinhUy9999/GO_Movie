package models

import (
	"my-app/config"
)

type Movie struct {
	MovieID  int    `json:"movie_id"`
	Title    string `json:"title"`
	Genre    string `json:"genre"`
	Duration int    `json:"duration"` // in minutes
	Picture  string `json:"picture"`
}

func CreateMovie(movie Movie) error {
	_, err := config.DB.Exec(
		"INSERT INTO MOVIE (title, genre, duration, picture) VALUES (?, ?, ?, ?)",
		movie.Title, movie.Genre, movie.Duration, movie.Picture,
	)
	if err != nil {
		return err
	}
	return nil
}
func GetAllMovies() ([]Movie, error) {
	rows, err := config.DB.Query("SELECT movieID, title, genre, duration, picture FROM MOVIE")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var movie Movie
		err := rows.Scan(&movie.MovieID, &movie.Title, &movie.Genre, &movie.Duration, &movie.Picture)
		if err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}
func UpdateMovie(movie Movie) error {
	_, err := config.DB.Exec(
		"UPDATE MOVIE SET title = ?, genre = ?, duration = ?, picture = ? WHERE movieID = ?",
		movie.Title, movie.Genre, movie.Duration, movie.Picture, movie.MovieID,
	)
	if err != nil {
		return err
	}
	return nil
}
func DeleteMovie(movieID int) error {
	_, err := config.DB.Exec("DELETE FROM MOVIE WHERE movieID = ?", movieID)
	if err != nil {
		return err
	}
	return nil
}
