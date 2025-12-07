package models

import (
	"errors"
	"time"

	"github.com/NikosGour/chatter/internal/common"
	"github.com/google/uuid"
)

var (
	ErrMessageNotFound = errors.New("message not found")
)

type Message struct {
	Id        int64     `json:"id,omitempty"`
	Sender    uuid.UUID `validate:"required" json:"sender_id,omitempty"`
	Recipient uuid.UUID `validate:"required" json:"recipient_id,omitempty"`
	DateSent  time.Time `validate:"required" json:"date_sent,omitempty"`
}

func (m Message) Validate() error {
	err := common.Validate.Struct(m)
	if err != nil {
		return err
	}
	return nil
}
