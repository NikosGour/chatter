package services

import (
	"database/sql"
	"errors"

	"github.com/NikosGour/chatter/internal/models"
	"github.com/NikosGour/chatter/internal/repositories"
	"github.com/NikosGour/logging/log"
	"github.com/google/uuid"
)

type GroupService struct {
	group_repo repositories.GroupRepository

	uuid_generator *ChannelService
	user_service   *UserService
}

func NewGroupService(group_repo repositories.GroupRepository, uuid_generator *ChannelService, user_service *UserService) *GroupService {
	s := &GroupService{group_repo: group_repo, uuid_generator: uuid_generator, user_service: user_service}
	return s
}

// Retrieves all group records from the database.
//
// Might return any sql error.
func (s *GroupService) GetAll() ([]models.Group, error) {
	group_dbos, err := s.group_repo.GetAll()
	if err != nil {
		return nil, err
	}

	groups := []models.Group{}
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
func (s *GroupService) GetByID(id uuid.UUID) (*models.Group, error) {
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
func (s *GroupService) Create(group *models.Group) (uuid.UUID, error) {
	id, err := s.uuid_generator.Create(models.ChannelTypeGroup)
	if err != nil {
		return uuid.Nil, err
	}
	group.Id = id
	gdbo := groupToDBO(group)
	return s.group_repo.Create(gdbo)

}

// Adds a the user of the given UUID to the list of subscribed users of the group
//
// Might return ErrGroupNotFound or any other sql error
func (s *GroupService) AddUserToGroup(user_id uuid.UUID, group_id uuid.UUID) error {
	_, err := s.user_service.GetByID(user_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ErrUserNotFound
		}
		return err
	}

	_, err = s.GetByID(group_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ErrGroupNotFound
		}
		return err
	}

	return s.group_repo.AddUserToGroup(user_id, group_id)
}

// Get all the user UUIDs from a group's user list
//
// Might return ErrGroupHasNoUsers or any other sql error
func (s *GroupService) GetUsers(group_id uuid.UUID) ([]models.User, error) {
	user_ids, err := s.group_repo.GetUsers(group_id)
	if err != nil {
		return nil, err
	}

	users := []models.User{}
	for _, user_id := range user_ids {
		user_dbo, err := s.user_service.GetByID(user_id)
		if err != nil {
			if errors.Is(err, models.ErrUserNotFound) {
				log.Warn("while getting users for group: %s, tried to get missing user: %s", group_id, user_id)
			}
			return nil, err
		}
		user := s.user_service.ToUser(user_dbo)
		users = append(users, *user)
	}
	return users, nil
}

// Transforms a group DBO to a group model
func (s *GroupService) toGroup(group_dbo repositories.GroupDBO) (*models.Group, error) {
	group := &group_dbo
	users, err := s.GetUsers(group.Id)
	if err != nil {
		return nil, err
	}
	group.Users = users
	return group, nil
}

func groupToDBO(g *models.Group) *repositories.GroupDBO {
	return g
}
