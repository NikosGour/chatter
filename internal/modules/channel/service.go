package channel

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type Service struct {
	channel_repo Repository
}

func NewService(channel_repo Repository) *Service {
	s := &Service{channel_repo: channel_repo}
	return s
}

func (s *Service) Create(chtype ChannelType) (uuid.UUID, error) {
	id, err := s.createNewUUID()
	if err != nil {
		return uuid.Nil, fmt.Errorf("On createNewUUID: %w", err)
	}

	return s.channel_repo.Create(id, chtype)
}

func (s *Service) createNewUUID() (uuid.UUID, error) {
	for {
		id := uuid.New()
		u, err := s.channel_repo.GetByID(id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return id, nil
			}
			return uuid.Nil, fmt.Errorf("On GetById: %w", err)
		}
		if u == nil {
			return id, nil
		}
	}
}
