package models

import (
	"bytes"
	"encoding/json"
	"fmt"

	jsonv2 "github.com/go-json-experiment/json"
	"github.com/go-playground/validator/v10"
)

type User struct {
	Username string `json:"login" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
}

func InitNewUser(data []byte) (*User, error) {
	var buf bytes.Buffer
	buf.Write(data)
	decoder := json.NewDecoder(&buf)
	decoder.DisallowUnknownFields()

	u := &User{}
	err := decoder.Decode(&u)
	if err != nil {
		return nil, err
	}

	validate := validator.New()
	err = validate.Struct(u)
	if err != nil {
		return nil, err

	}

	err = checkDuplicateFields(data)
	if err != nil {
		return nil, err

	}
	return u, nil
}

func checkDuplicateFields(data []byte) error {
	u := &User{}
	fmt.Println(string(data))
	err := jsonv2.UnmarshalOptions{}.Unmarshal(jsonv2.DecodeOptions{
		AllowDuplicateNames: false,
	}, data, &u)

	if err != nil {
		return err
	}
	return nil
}
