package models

import (
	"errors"
	"time"

	"github.com/NikosGour/chatter/internal/common"
)

var (
	ErrMessageNotFound = errors.New("message not found")
)

type Message struct {
	Id       int64     `json:"id,omitempty"`
	Text     string    `json:"text"`
	Sender   *User     `json:"sender,omitempty"`
	Tab      *Tab      `json:"tab,omitempty"`
	DateSent time.Time `validate:"required" json:"date_sent,omitempty,omitzero"`
}

func (m Message) Validate() error {
	err := common.Validate.Struct(m)
	if err != nil {
		return err
	}
	return nil
}
