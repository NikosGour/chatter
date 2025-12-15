package services

import (
	"errors"
	"fmt"

	"github.com/NikosGour/chatter/internal/models"
	"github.com/NikosGour/chatter/internal/repositories"
	"github.com/NikosGour/logging/log"
	"github.com/google/uuid"
)

type ServerService struct {
	server_repo repositories.ServerRepository

	user_service *UserService
	tab_service  *TabService
}

func NewServerService(server_repo repositories.ServerRepository, user_service *UserService, tab_service *TabService) *ServerService {
	s := &ServerService{server_repo: server_repo, user_service: user_service, tab_service: tab_service}
	return s
}

// Retrieves all servers from the database.
//
// Might return any sql error.
func (s *ServerService) GetAll() ([]models.Server, error) {
	server_dbos, err := s.server_repo.GetAll()
	if err != nil {
		return nil, err
	}

	servers := []models.Server{}
	for _, server_dbo := range server_dbos {
		server, err := s.toServer(server_dbo)
		if err != nil {
			return nil, err
		}
		servers = append(servers, *server)
	}
	return servers, nil
}

// Retrieves a server given the UUID.
//
// Might return ErrServerNotFound or any other sql error
func (s *ServerService) GetByID(id uuid.UUID) (*models.Server, error) {
	Server_dbo, err := s.server_repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	Server, err := s.toServer(*Server_dbo)
	if err != nil {
		return nil, err
	}
	return Server, nil
}

// Inserts a server into a database.
//
// Returns the UUID of the created server.
// Might return any sql error
func (s *ServerService) Create(server *models.Server) (uuid.UUID, error) {
	id, err := s.generateUUID()
	if err != nil {
		return uuid.Nil, err
	}
	server.Id = id
	server_dbo := ServerToDBO(server)
	return s.server_repo.Create(server_dbo)

}

// Adds a the user of the given UUID to the list of subscribed users of the server
//
// Might return ErrServerNotFound or any other sql error
func (s *ServerService) AddUserToServer(user_id uuid.UUID, server_id uuid.UUID) error {
	_, err := s.user_service.GetByID(user_id)
	if err != nil {
		return err
	}

	_, err = s.GetByID(server_id)
	if err != nil {
		return err
	}

	return s.server_repo.AddUserToServer(user_id, server_id)
}

// Get all the user UUIDs from a Server's user list
//
// Might return ErrServerHasNoUsers or any other sql error
func (s *ServerService) GetUsers(server_id uuid.UUID) ([]models.User, error) {
	user_ids, err := s.server_repo.GetUsers(server_id)
	if err != nil {
		return nil, err
	}

	users := []models.User{}
	for _, user_id := range user_ids {
		user_dbo, err := s.user_service.GetByID(user_id)
		if err != nil {
			if errors.Is(err, models.ErrUserNotFound) {
				log.Warn("while getting users for Server: %s, tried to get missing user: %s", server_id, user_id)
			}
			return nil, err
		}
		user := s.user_service.ToUser(user_dbo)
		users = append(users, *user)
	}
	return users, nil
}

// Get all the user UUIDs from a Server's user list
//
// Might return ErrServerHasNoUsers or any other sql error
func (s *ServerService) GetTabs(server_id uuid.UUID) ([]models.Tab, error) {
	tabs, err := s.tab_service.GetByServerID(server_id)
	if err != nil {
		return nil, err
	}

	return tabs, nil
}

// Transforms a Server DBO to a server model
func (s *ServerService) toServer(server_dbo repositories.ServerDBO) (*models.Server, error) {
	server := &server_dbo
	users, err := s.GetUsers(server.Id)
	if err != nil {
		return nil, err
	}
	server.Users = users
	return server, nil
}

func ServerToDBO(s *models.Server) *repositories.ServerDBO {
	return s
}

func (s *ServerService) generateUUID() (uuid.UUID, error) {
	id := uuid.New()

	for {
		ch, err := s.server_repo.GetByID(id)
		if err != nil {
			if errors.Is(err, models.ErrServerNotFound) {
				break
			}
			return uuid.Nil, fmt.Errorf("On GetById: %w", err)
		}

		if ch == nil {
			break
		}
		id = uuid.New()
	}

	return id, nil
}
