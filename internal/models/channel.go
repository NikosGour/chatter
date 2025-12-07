package models

import "github.com/google/uuid"

type Channel struct {
	Id          uuid.UUID   `db:"id"`
	ChannelType ChannelType `db:"channel_type"`
}

type ChannelType string

const (
	ChannelTypeUser  ChannelType = "user"
	ChannelTypeGroup ChannelType = "group"
)
