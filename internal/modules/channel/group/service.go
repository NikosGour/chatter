package group

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/NikosGour/chatter/internal/modules/channel"
	"github.com/NikosGour/chatter/internal/modules/channel/user"
	"github.com/NikosGour/logging/log"
	"github.com/google/uuid"
)

type Service struct {
	group_repo   Repository
	channel_repo channel.Repository
	user_repo    user.Repository
}

func NewService(group_repo Repository, channel_repo channel.Repository, user_repo user.Repository) *Service {
	s := &Service{group_repo: group_repo, channel_repo: channel_repo, user_repo: user_repo}
	return s
}

// Retrieves all group records from the database.
//
// Might return any sql error.
func (s *Service) GetAll() ([]Group, error) {
	group_dbos, err := s.group_repo.GetAll()
	if err != nil {
		return nil, err
	}

	groups := []Group{}
	for _, group_dbo := range group_dbos {
		group, err := s.toGroup(group_dbo)
		if err != nil {
			return nil, err
		}
		groups = append(groups, *group)
	}
	return groups, nil
}

// Retrieves a group given the UUID.
//
// Might return ErrGroupNotFound or any other sql error
func (s *Service) GetByID(id uuid.UUID) (*Group, error) {
	group_dbo, err := s.group_repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	group, err := s.toGroup(*group_dbo)
	if err != nil {
		return nil, err
	}
	return group, nil
}

// Inserts a group into a database.
//
// Returns the UUID of the created group.
// Might return any sql error
func (s *Service) Create(group *Group) (uuid.UUID, error) {
	id, err := s.channel_repo.Create(channel.ChannelTypeGroup)
	if err != nil {
		return uuid.Nil, fmt.Errorf("On channel create: %w", err)
	}
	group.Id = id
	return s.group_repo.Create(group)

}

// Adds a the user of the given UUID to the list of subscribed users of the group
//
// Might return ErrGroupNotFound or any other sql error
func (s *Service) AddUserToGroup(user_id uuid.UUID, group_id uuid.UUID) error {
	_, err := s.user_repo.GetByID(user_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user.ErrUserNotFound
		}
		return err
	}

	_, err = s.GetByID(group_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrGroupNotFound
		}
		return err
	}

	return s.group_repo.AddUserToGroup(user_id, group_id)
}

// Get all the user UUIDs from a group's user list
//
// Might return ErrGroupHasNoUsers or any other sql error
func (s *Service) GetUsers(group_id uuid.UUID) ([]user.User, error) {
	user_ids, err := s.group_repo.GetUsers(group_id)
	if err != nil {
		return nil, err
	}

	users := []user.User{}
	for _, user_id := range user_ids {
		user_dbo, err := s.user_repo.GetByID(user_id)
		if err != nil {
			if errors.Is(err, user.ErrUserNotFound) {
				log.Warn("while getting users for group: %s, tried to get missing user: %s", group_id, user_id)
			}
			return nil, err
		}
		user := s.user_repo.ToUser(user_dbo)
		users = append(users, *user)
	}
	return users, nil
}

// Transforms a group DBO to a group model
func (s *Service) toGroup(group_dbo groupDBO) (*Group, error) {
	group := &group_dbo
	users, err := s.GetUsers(group.Id)
	if err != nil {
		return nil, err
	}
	group.Users = users
	return group, nil
}

func (g *Group) toDBO() *groupDBO {
	return g
}
