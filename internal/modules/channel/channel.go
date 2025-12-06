package channel

import (
	"github.com/google/uuid"
)

type Channel interface {
	GetId() uuid.UUID
}

type ChannelType string

const (
	ChannelTypeUser  ChannelType = "user"
	ChannelTypeGroup ChannelType = "group"
)
