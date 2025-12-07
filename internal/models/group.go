package models

import (
	"errors"
	"time"

	"github.com/NikosGour/chatter/internal/common"
	"github.com/google/uuid"
)

var (
	ErrGroupNotFound   = errors.New("group not found")
	ErrGroupHasNoUsers = errors.New("group has no users")
)

type Group struct {
	Id          uuid.UUID `json:"id,omitempty" db:"id"`
	Name        string    `validate:"required" json:"name,omitempty" db:"name"`
	Users       []User    `json:"users,omitempty" db:"users"`
	DateCreated time.Time `validate:"required" json:"date_created,omitempty" db:"date_created"`
}

func (g *Group) GetId() uuid.UUID {
	return g.Id
}

func (g Group) Validate() error {
	err := common.Validate.Struct(g)
	if err != nil {
		return err
	}
	return nil
}
