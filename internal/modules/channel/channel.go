package channel

import (
	"github.com/google/uuid"
)

type Channel interface {
	GetId() uuid.UUID
}
