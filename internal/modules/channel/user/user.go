package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id          uuid.UUID
	Username    string
	Password    string
	DateCreated time.Time
}

func (u *User) GetId() uuid.UUID {
	return u.Id
}
