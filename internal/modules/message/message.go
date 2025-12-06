package message

import (
	"errors"
	"time"

	"github.com/NikosGour/chatter/internal/modules/channel"
)

var (
	ErrMessageNotFound = errors.New("message not found")
)

type Message struct {
	Id        int64
	Sender    channel.Channel
	Recipient channel.Channel
	DateSent  time.Time
}
