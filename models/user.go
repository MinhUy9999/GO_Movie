package models

import (
	"database/sql"
	"errors"
	"my-app/config"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`   // Either "user" or "admin"
	Gender   string `json:"gender"` // Male, Female, Other
}

// Register a new user
func RegisterUser(user User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return err
	}

	_, err = config.DB.Exec("INSERT INTO users (email, name, password, phone, role, gender) VALUES (?, ?, ?, ?, ?, ?)",
		user.Email, user.Name, hashedPassword, user.Phone, user.Role, user.Gender)
	return err
}

// Authenticate a user (login)
func AuthenticateUser(email, password string) (User, bool, error) {
	var user User
	err := config.DB.QueryRow("SELECT id, email, password, role FROM users WHERE email = ?", email).Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if err == sql.ErrNoRows {
		return user, false, errors.New("user not found")
	} else if err != nil {
		return user, false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, false, errors.New("incorrect password")
	}

	return user, true, nil
}

// GetUserByID retrieves a user by their ID
func GetUserByID(id int) (User, error) {
	var user User
	err := config.DB.QueryRow(
		"SELECT id, email, name, phone, gender, role FROM users WHERE id = ?", id,
	).Scan(&user.ID, &user.Email, &user.Name, &user.Phone, &user.Gender, &user.Role)

	if err == sql.ErrNoRows {
		return user, errors.New("user not found")
	} else if err != nil {
		return user, err
	}
	return user, nil
}

// DeleteUser deletes a user by ID
func DeleteUser(userID int) error {
	_, err := config.DB.Exec("DELETE FROM users WHERE id = ?", userID)
	return err
}

// UpdateUser updates user details in the database
func UpdateUser(user User) error {
	_, err := config.DB.Exec("UPDATE users SET name = ?, email = ?, phone = ? WHERE id = ?",
		user.Name, user.Email, user.Phone, user.ID)
	return err
}

// GetAllUsers retrieves all users from the database
func GetAllUsers() ([]User, error) {
	rows, err := config.DB.Query("SELECT id, email, name, phone, gender, role FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.Phone, &user.Gender, &user.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
