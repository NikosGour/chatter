package models

import "github.com/google/uuid"

type Channel interface {
	GetId() uuid.UUID
	GetName() string
}
type ChannelType string

const (
	ChannelTypeUser  ChannelType = "user"
	ChannelTypeGroup ChannelType = "group"
)
