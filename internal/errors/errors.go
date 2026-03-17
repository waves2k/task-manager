package errors

import (
	"fmt"
)

type ErrorType string

const (
	TypeNotFound         ErrorType = "NOT_FOUND"
	TypeAlreadyExists    ErrorType = "ALREADY_EXISTS"
	TypeValidation       ErrorType = "VALIDATION"
	TypeUnauthorized     ErrorType = "UNAUTHORIZED"
	TypeForbidden        ErrorType = "FORBIDDEN"
	TypeInternal         ErrorType = "INTERNAL"
	TypeWrongInputData   ErrorType = "INVALID_INPUT_DATA"
	TypeWrongCredentials ErrorType = "INVALID_CREDENTIALS"
)

type AppError struct {
	Type    ErrorType
	Message string
	Err     error
}

func (e AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s : %e", e.Message, e.Err)
	}
	return e.Message
}

func (e AppError) Unwrap() error {
	return e.Err
}

func NewWrongInputDataError(message string, err error) AppError {
	return AppError{
		Type:    TypeWrongInputData,
		Message: message,
		Err:     err,
	}
}

func NewNotFoundError(message string, err error) AppError {
	return AppError{
		Type:    TypeNotFound,
		Message: message,
		Err:     err,
	}
}

func NewAlreadyExistsError(message string, err error) AppError {
	return AppError{
		Type:    TypeAlreadyExists,
		Message: message,
		Err:     err,
	}
}

func NewValidationError(message string, err error) AppError {
	return AppError{
		Type:    TypeValidation,
		Message: message,
		Err:     err,
	}
}

func NewUnauthorizedError(message string, err error) AppError {
	return AppError{
		Type:    TypeUnauthorized,
		Message: message,
		Err:     err,
	}
}

func NewForbiddenError(message string, err error) AppError {
	return AppError{
		Type:    TypeForbidden,
		Message: message,
		Err:     err,
	}
}

func NewInternalError(message string, err error) AppError {
	return AppError{
		Type:    TypeInternal,
		Message: message,
		Err:     err,
	}
}
