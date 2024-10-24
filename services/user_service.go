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

// ValidateUsernamePassword проверяет корректность username и password
// Allowed characters: a-z, A-Z, 0-9, ., _, -.
func ValidateUsernamePassword(username, password string) error {
	// Проверка длины username
	if len(username) < 5 || len(username) > 30 {
		return ErrInvalidUsername
	}
	// Проверка допустимых символов в username
	validUsername := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	if !validUsername.MatchString(username) {
		return ErrInvalidUsername
	}

	// Проверка длины password
	if len(password) < 5 || len(password) > 30 {
		return ErrInvalidPassword
	}
	// Проверка допустимых символов в password (можно адаптировать под требования)
	validPassword := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	if !validPassword.MatchString(password) {
		return ErrInvalidPassword
	}

	return nil
}
