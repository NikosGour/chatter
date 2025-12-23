package models

import (
	"errors"
	"time"

	"github.com/NikosGour/chatter/internal/common"
	"github.com/google/uuid"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type User struct {
	Id          uuid.UUID `json:"id,omitempty" db:"id"`
	Username    string    `json:"username,omitempty" db:"username"`
	Password    string    `json:"password,omitempty" db:"password"`
	DateCreated time.Time `json:"date_created,omitempty,omitzero" db:"date_created"`
	IsTest      bool      `db:"is_test"`
}

func (u User) Validate() error {
	err := common.Validate.Struct(u)
	if err != nil {
		return err
	}
	return nil
}
