package group

import (
	"time"

	"github.com/NikosGour/chatter/internal/modules/channel/user"
	"github.com/google/uuid"
)

type Group struct {
	Id          uuid.UUID
	Name        string
	Users       []user.User
	DateCreated time.Time
}

func (g *Group) GetId() uuid.UUID {
	return g.Id
}
