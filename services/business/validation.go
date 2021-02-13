// Package business implements different business services required by the user service
package business

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// Validate validates the CreateUserRequest model and return error if the validation failes
// Returns error if validation failes
func (val CreateUserRequest) Validate() error {
	return validation.ValidateStruct(&val,
		// Validate User using its own validation rules
		validation.Field(&val.User),
	)
}

// Validate validates the ReadUserRequest model and return error if the validation failes
// Returns error if validation failes
func (val ReadUserRequest) Validate() error {
	return validation.ValidateStruct(&val,
		// UserID cannot be empty
		validation.Field(&val.UserID, validation.Required),
	)
}

// Validate validates the ReadUserByEmailRequest model and return error if the validation failes
// Returns error if validation failes
func (val ReadUserByEmailRequest) Validate() error {
	return validation.ValidateStruct(&val,
		// Email address cannot be empty
		validation.Field(&val.Email, validation.Required, is.Email),
	)
}

// Validate validates the UpdateUserRequest model and return error if the validation failes
// Returns error if validation failes
func (val UpdateUserRequest) Validate() error {
	return validation.ValidateStruct(&val,
		// UserID cannot be empty
		validation.Field(&val.UserID, validation.Required),
		// Validate User using its own validation rules
		validation.Field(&val.User),
	)
}

// Validate validates the DeleteUserRequest model and return error if the validation failes
// Returns error if validation failes
func (val DeleteUserRequest) Validate() error {
	return validation.ValidateStruct(&val,
		// UserID cannot be empty
		validation.Field(&val.UserID, validation.Required),
	)
}

// Validate validates the SearchRequest model and return error if the validation failes
// Returns error if validation failes
func (val SearchRequest) Validate() error {
	return nil
}
