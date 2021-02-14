// Package models defines the different object models used in User
package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

// Validate validates the User and return error if the validation failes
// Returns error if validation failes
func (val User) Validate() error {
	return validation.ValidateStruct(&val)
}
