package models

import (
	"errors"
)

type ResponseBalanceError struct {
	Text string
	Err  error
}

func (e *ResponseBalanceError) Error() string { return e.Text + ": " + e.Err.Error() }

func ErrNoBalance() error {
	return &ResponseBalanceError{
		Text: "the user has no balance",
		Err:  errors.New("no balance"),
	}
}

type ResponseOrderError struct {
	Text string
	Err  error
}

func (e *ResponseOrderError) Error() string { return e.Text + ": " + e.Err.Error() }

func ErrNoOrder() error {
	return &ResponseOrderError{
		Text: "the user has no balance",
		Err:  errors.New("no balance"),
	}
}
