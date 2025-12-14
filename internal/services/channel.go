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

	user_service  *UserService
	group_service *GroupService
}

func NewChannelService(channel_repo repositories.ChannelRepository) *ChannelService {
	return &ChannelService{channel_repo: channel_repo}
}
func (s *ChannelService) AddUserService(u *UserService) {
	s.user_service = u
}

func (s *ChannelService) AddGroupService(g *GroupService) {
	s.group_service = g
}

func (s *ChannelService) Create(chtype models.ChannelType) (uuid.UUID, error) {
	id := uuid.New()

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

func (s *ChannelService) GetByID(id uuid.UUID) (models.Channel, error) {
	chdbo, err := s.channel_repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.toChannel(chdbo)
}

func (s *ChannelService) toChannel(chdbo *repositories.ChannelDBO) (models.Channel, error) {
	switch chdbo.ChannelType {
	case models.ChannelTypeUser:
		user, err := s.user_service.GetByID(chdbo.Id)
		if err != nil {
			return nil, err
		}

		return user, nil

	case models.ChannelTypeGroup:
		group, err := s.group_service.GetByID(chdbo.Id)

		if err != nil {
			return nil, err
		}

		return group, nil
	default:
		return nil, fmt.Errorf("%s:%s", models.ErrInvalidChannelType, chdbo.ChannelType)
	}

}
