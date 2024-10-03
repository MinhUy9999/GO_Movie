package utils // Ensure package declaration

import (
	"fmt"
	"strconv"
)

// Example helper function to convert a string to an int
func StringToInt(s string) (int, error) {
	value, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("failed to convert string to int: %v", err)
	}
	return value, nil
}

// Example helper function to check for errors
func CheckError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
	}
}
