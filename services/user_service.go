package services

import (
	"errors"
	"manzi/repositories"
	"regexp"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidUsername   = errors.New("invalid username")
	ErrInvalidPassword   = errors.New("invalid password")
)

// RegisterUser handles the logic of checking and registering a new user.
func RegisterUser(username, password string) error {
	exists, err := repositories.CheckUserExists(username)
	if err != nil {
		return err
	}

	if exists {
		return ErrUserAlreadyExists
	}

	err = repositories.CreateUser(username, password)
	if err != nil {
		return err
	}

	return nil
}

// ValidateUsernamePassword checks the validity of username and password.
// Allowed characters: a-z, A-Z, 0-9, ., _, -.
func ValidateUsernamePassword(username, password string) error {
	// Check username length
	if len(username) < 5 || len(username) > 30 {
		return ErrInvalidUsername
	}
	// Check valid characters in username
	validUsername := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	if !validUsername.MatchString(username) {
		return ErrInvalidUsername
	}

	// Check password length
	if len(password) < 5 || len(password) > 30 {
		return ErrInvalidPassword
	}
	// Check valid characters in password (can be adapted as needed)
	validPassword := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	if !validPassword.MatchString(password) {
		return ErrInvalidPassword
	}

	return nil
}
