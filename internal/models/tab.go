package models

import (
	"errors"
	"time"

	"github.com/NikosGour/chatter/internal/common"
	"github.com/google/uuid"
)

var (
	ErrTabNotFound = errors.New("tab not found")
)

type Tab struct {
	Id          uuid.UUID `json:"id,omitempty" db:"id"`
	Name        string    `json:"name,omitempty" db:"name"`
	ServerId    uuid.UUID `json:"server_id,omitempty" db:"server_id"`
	Server      *Server   `json:"server,omitempty" db:"server"`
	DateCreated time.Time `json:"date_created,omitempty,omitzero" db:"date_created"`
}

func (t Tab) Validate() error {
	err := common.Validate.Struct(t)
	if err != nil {
		return err
	}
	return nil
}
