package models

import (
	"errors"
	"time"

	"github.com/NikosGour/chatter/internal/common"
	"github.com/google/uuid"
)

var (
	ErrServerNotFound   = errors.New("server not found")
	ErrServerHasNoUsers = errors.New("server has no users")
)

type Server struct {
	Id          uuid.UUID `json:"id,omitempty" db:"id"`
	Name        string    `json:"name,omitempty" db:"name"`
	Users       []User    `json:"users,omitempty" db:"users"`
	Tabs        []Tab     `json:"tabs,omitempty" db:"tabs"`
	DateCreated time.Time `json:"date_created,omitempty" db:"date_created"`
}

func (s Server) Validate() error {
	err := common.Validate.Struct(s)
	if err != nil {
		return err
	}
	return nil
}
