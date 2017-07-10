package gocqrs

import (
	"errors"
	"github.com/diegogub/trk/lib"
)

const (
	AccEntity = "acc"
)

type Account struct {
	ID      string                 `json:"id"`
	Name    string                 `json:"name"`
	Active  bool                   `json:"active"`
	AdminID string                 `json:"admid"`
	Data    map[string]interface{} `json:"data"`
}

func (acc *Account) GetName() string {
	return "acc_validator"
}

func (acc *Account) Validate(e Entity) error {
	e.Decode(acc)

	if acc.Name == "" {
		return errors.New("Invalid account name")
	}

	err := lib.ValidID(acc.AdminID, true)
	if err != nil {
		return err
	}

	return err
}
