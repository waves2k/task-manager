package errors

import "errors"

func IsWrongInputData(err error) bool {
	var appErr AppError
	return errors.As(err, &appErr) && appErr.Type == TypeWrongInputData
}

func IsNotFound(err error) bool {
	var appErr AppError
	return errors.As(err, &appErr) && (appErr.Type == UserNotFoundMessage ||
		appErr.Type == TodoNotFoundMessage ||
		appErr.Type == ListNotFoundMessage)
}

func IsAlreadyExists(err error) bool {
	var appErr AppError
	return errors.As(err, &appErr) && (appErr.Type == UserAlreadyExistsMessage ||
		appErr.Type == TodoAlreadyExistsMessage ||
		appErr.Type == ListAlreadyExistsMessage)
}

func IsValidation(err error) bool {
	var appErr AppError
	return errors.As(err, &appErr) && (appErr.Type == TypeValidation)
}

func IsUnauthorized(err error) bool {
	var appErr AppError
	return errors.As(err, &appErr) && (appErr.Type == TypeUnauthorized)
}

func IsFrobidden(err error) bool {
	var appErr AppError
	return errors.As(err, &appErr) && (appErr.Type == TypeForbidden)
}

func IsInternal(err error) bool {
	var appErr AppError
	return errors.As(err, &appErr) && (appErr.Type == TypeInternal)
}
