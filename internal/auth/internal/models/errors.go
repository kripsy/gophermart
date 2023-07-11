package models

import "fmt"

type UserExistsError struct {
	Text string
	Err  error
}

func NewUserExistsError(username string, err error) error {
	return &UserExistsError{
		Text: fmt.Sprintf("%v already exists", username),
		Err:  err,
	}
}

func (ue *UserExistsError) Error() string {
	return ue.Err.Error()
}
