package models

import (
	"encoding/json"
	"fmt"
)

type User struct {
	Username string `json:"login"`
	Password string `json:"password,omitempty"`
}

func InitNewUser(data []byte) (*User, error) {
	u := &User{}
	err := json.Unmarshal(data, u)
	if err != nil {
		return nil, err
	}

	if u.Password == "" {
		return nil, fmt.Errorf("user password is empty")
	}
	return u, nil
}
