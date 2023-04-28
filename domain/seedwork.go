package domain

import "time"

// TODO: is idiomatic naming ErrValidation or ValidationErr?
type ValidationError struct {
	Field   string
	Message string
}

func NewValidationError(field string, message string) ValidationError {
	return ValidationError{
		Field:   field,
		Message: message,
	}
}

func (e ValidationError) Error() string {
	return "validation error string"
}

type Entity struct {
	// TODO: add Id of generic type?
	CreatedAt time.Time
	UpdatedAt *time.Time
}

// TODO: add equals method to entity, comparing Ids.

type AggregateRoot struct {
	Entity
	Events []any
}
