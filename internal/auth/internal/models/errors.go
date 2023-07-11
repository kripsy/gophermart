package models

type UserExistsError struct {
	Text string
	Err  error
}

func NewUserExistsError(username string, err error) {

}
