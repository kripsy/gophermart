package models

import (
	"errors"
	"fmt"
)

type UserExistsError struct {
	Text string
	Err  error
}

func NewUserExistsError(username string) error {
	return &UserExistsError{
		Text: fmt.Sprintf("%v already exists", username),
		Err:  errors.New("user already exists"),
	}
}

func (ue *UserExistsError) Error() string {
	return ue.Err.Error()
}

type UserLoginError struct {
	Text string
	Err  error
}

func NewUserLoginError(username string) error {
	return &UserLoginError{
		Text: fmt.Sprintf("login failed for %v", username),
		Err:  errors.New("login failed"),
	}
}

func (ue *UserLoginError) Error() string {
	return ue.Err.Error()
}
