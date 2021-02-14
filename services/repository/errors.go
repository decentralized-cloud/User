// Package repository implements different repository services required by the user service
package repository

import "fmt"

// UnknownError indicates that an unknown error has happened
type UnknownError struct {
	Message string
	Err     error
}

// Error returns message for the UnknownError error type
// Returns the error nessage
func (e UnknownError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("Unknown error occurred. Error message: %s.", e.Message)
	}

	return fmt.Sprintf("Unknown error occurred. Error message: %s. Error: %s", e.Message, e.Err.Error())
}

// Unwrap returns the err if provided through NewUnknownErrorWithError function, otherwise returns nil
func (e UnknownError) Unwrap() error {
	return e.Err
}

// IsUnknownError indicates whether the error is of type UnknownError
func IsUnknownError(err error) bool {
	_, ok := err.(UnknownError)

	return ok
}

// NewUnknownError creates a new UnknownError error
func NewUnknownError(message string) error {
	return UnknownError{
		Message: message,
	}
}

// NewUnknownErrorWithError creates a new UnknownError error
func NewUnknownErrorWithError(message string, err error) error {
	return UnknownError{
		Message: message,
		Err:     err,
	}
}

// UserAlreadyExistsError indicates that the user with the given information already exists
type UserAlreadyExistsError struct {
	Err error
}

// Error returns message for the UserAlreadyExistsError error type
// Returns the error nessage
func (e UserAlreadyExistsError) Error() string {
	if e.Err == nil {
		return "User already exists."
	}

	return fmt.Sprintf("User already exists. Error: %s", e.Err.Error())
}

// Unwrap returns the err if provided through NewUserAlreadyExistsErrorWithError function, otherwise returns nil
func (e UserAlreadyExistsError) Unwrap() error {
	return e.Err
}

// IsUserAlreadyExistsError indicates whether the error is of type UserAlreadyExistsError
func IsUserAlreadyExistsError(err error) bool {
	_, ok := err.(UserAlreadyExistsError)

	return ok
}

// NewUserAlreadyExistsError creates a new UserAlreadyExistsError error
func NewUserAlreadyExistsError() error {
	return UserAlreadyExistsError{}
}

// NewUserAlreadyExistsErrorWithError creates a new UserAlreadyExistsError error
func NewUserAlreadyExistsErrorWithError(err error) error {
	return UserAlreadyExistsError{
		Err: err,
	}
}

// UserNotFoundError indicates that the user with the given email address does not exist
type UserNotFoundError struct {
	Email string
	Err   error
}

// Error returns message for the UserNotFoundError error type
// Returns the error nessage
func (e UserNotFoundError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("User not found. Email: %s.", e.Email)
	}

	return fmt.Sprintf("User not found. Email: %s. Error: %s", e.Email, e.Err.Error())
}

// Unwrap returns the err if provided through UserNotFoundError function, otherwise returns nil
func (e UserNotFoundError) Unwrap() error {
	return e.Err
}

// IsUserNotFoundError indicates whether the error is of type UserNotFoundError
func IsUserNotFoundError(err error) bool {
	_, ok := err.(UserNotFoundError)

	return ok
}

// NewUserNotFoundError creates a new UserNotFoundError error
// email: Mandatory. The email address that did not match any existing user
func NewUserNotFoundError(email string) error {
	return UserNotFoundError{
		Email: email,
	}
}

// NewUserNotFoundErrorWithError creates a new UserNotFoundError error
// email: Mandatory. The email address that did not match any existing user
func NewUserNotFoundErrorWithError(email string, err error) error {
	return UserNotFoundError{
		Email: email,
		Err:   err,
	}
}
