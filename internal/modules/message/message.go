package message

import (
	"time"

	"github.com/NikosGour/chatter/internal/modules/channel"
)

type Message struct {
	Id        int
	Sender    channel.Channel
	Recipient channel.Channel
	DateSent  time.Time
}
