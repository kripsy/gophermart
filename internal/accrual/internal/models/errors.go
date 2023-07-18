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
		Text: "there are no accruals for the order",
		Err:  errors.New("there are no accruals for the order"),
	}
}
