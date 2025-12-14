package models

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrChannelNotFound = errors.New("channel not found")
)

type Channel interface {
	GetId() uuid.UUID
	GetName() string
	GetRecipients() []uuid.UUID
}

var (
	ErrInvalidChannelType = errors.New("invalid channel type")
)

type ChannelType string

const (
	ChannelTypeUser  ChannelType = "user"
	ChannelTypeGroup ChannelType = "group"
)
