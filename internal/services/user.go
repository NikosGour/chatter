package services

import (
	"github.com/NikosGour/chatter/internal/models"
	"github.com/NikosGour/chatter/internal/repositories"
	"github.com/google/uuid"
)

type User = models.User

type UserService struct {
	user_repo repositories.UserRepository

	uuid_generator *ChannelService
}

func NewUserService(user_repo repositories.UserRepository, uuid_generator *ChannelService) *UserService {
	s := &UserService{user_repo: user_repo, uuid_generator: uuid_generator}
	return s
}

// Operations
func (s *UserService) GetAll() ([]models.User, error) {
	udbos, err := s.user_repo.GetAll()
	if err != nil {
		return nil, err
	}

	us := []models.User{}
	for _, udbo := range udbos {
		u := s.ToUser(&udbo)
		us = append(us, *u)
	}

	return us, nil
}
func (s *UserService) GetByID(id uuid.UUID) (*models.User, error) {
	udbo, err := s.user_repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.ToUser(udbo), nil

}
func (s *UserService) Create(user *models.User) (uuid.UUID, error) {
	id, err := s.uuid_generator.Create(models.ChannelTypeUser)
	if err != nil {
		return uuid.Nil, err
	}
	user.Id = id

	udbo := userToDBO(user)
	return s.user_repo.Create(udbo)
}

func (s *UserService) ToUser(udb *repositories.UserDBO) *models.User {
	return udb
}
func userToDBO(u *models.User) *repositories.UserDBO {
	return u
}
