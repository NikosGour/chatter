package services

import (
	"errors"
	"fmt"

	"github.com/NikosGour/chatter/internal/models"
	"github.com/NikosGour/chatter/internal/repositories"
	"github.com/google/uuid"
)

type UserService struct {
	user_repo repositories.UserRepository
}

func NewUserService(user_repo repositories.UserRepository) *UserService {
	s := &UserService{user_repo: user_repo}
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
func (s *UserService) GetByUsername(username string) ([]models.User, error) {
	udbos, err := s.user_repo.GetByUsername(username)
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
func (s *UserService) GetByTestUsername(username string) ([]models.User, error) {
	udbos, err := s.user_repo.GetByTestUsername(username)
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

func (s *UserService) Create(user *models.User) (uuid.UUID, error) {
	id, err := s.generateUUID()
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

func (s *UserService) generateUUID() (uuid.UUID, error) {
	id := uuid.New()

	for {
		ch, err := s.user_repo.GetByID(id)
		if err != nil {
			if errors.Is(err, models.ErrUserNotFound) {
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
