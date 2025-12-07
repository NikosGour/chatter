package services

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/NikosGour/chatter/internal/models"
	"github.com/NikosGour/chatter/internal/repositories"
	"github.com/google/uuid"
)

type ChannelService struct {
	channel_repo repositories.ChannelRepository
}

func NewUUIDGenerator(channel_repo repositories.ChannelRepository) *ChannelService {
	return &ChannelService{channel_repo: channel_repo}
}

func (s *ChannelService) Create(chtype models.ChannelType) (uuid.UUID, error) {
	id, _ := uuid.Parse("be94832d-eb57-4214-a07b-647c5e64684f") //uuid.New()

	for {
		ch, err := s.channel_repo.GetByID(id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				break
			}
			return uuid.Nil, fmt.Errorf("On GetById: %w", err)
		}

		if ch == nil {
			break
		}
		id = uuid.New()
	}

	chdbo := repositories.ChannelDBO{
		Id:          id,
		ChannelType: chtype,
	}

	return s.channel_repo.Create(&chdbo)
}
