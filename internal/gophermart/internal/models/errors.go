package models

import (
	"errors"
	"fmt"
)

type ResponseBalanceError struct {
	Text string
	Err  error
}

func (e *ResponseBalanceError) Error() string { return e.Text + ": " + e.Err.Error() }

func ErrNoBalance() error {
	return &ResponseBalanceError{
		Text: fmt.Sprintf("the user has no balance"),
		Err:  errors.New("no balance"),
	}
}
