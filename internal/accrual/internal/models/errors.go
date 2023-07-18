package models

import (
	"errors"
)

type AccrualError struct {
	Text string
	Err  error
}

func (e *AccrualError) Error() string { return e.Text + ": " + e.Err.Error() }

func ErrNoAccrual() error {
	return &AccrualError{
		Text: "the user has no balance",
		Err:  errors.New("no balance"),
	}
}
