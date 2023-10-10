package model

import "errors"

var (
	ErrInvalidDateFormat = errors.New("invalid date format")
	ErrLoanNotFound      = errors.New("loan not found")
)
