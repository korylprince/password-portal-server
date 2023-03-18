package db

import "time"

// User represents a student user
type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

// SISDB represents a student information database
type SISDB interface {
	// GetID returns the id for the student matching the given information,
	// an empty string if the information doesn't match, or an error if one occurred
	GetID(firstName, lastName, ssn string, birthDate time.Time) (string, error)
}

// UserDB represents a user database
type UserDB interface {
	// Get returns the user with the given id, nil if the user doesn't exist, or an error if one occurred
	Get(id string) (*User, error)
}
