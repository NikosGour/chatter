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
	Username    string    `validate:"required" json:"username,omitempty" db:"username"`
	Password    string    `validate:"required" json:"password,omitempty" db:"password"`
	DateCreated time.Time `validate:"required" json:"date_created,omitempty" db:"date_created"`
}

func (u *User) GetId() uuid.UUID {
	return u.Id
}
func (u *User) GetName() string {
	return u.Username
}

func (u *User) GetRecipients() []uuid.UUID {
	return []uuid.UUID{u.Id}
}

func (u User) Validate() error {
	err := common.Validate.Struct(u)
	if err != nil {
		return err
	}
	return nil
}
