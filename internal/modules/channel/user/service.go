package user

import (
	"github.com/google/uuid"
)

type Service struct {
	user_repo Repository
}

func NewService(user_repo Repository) *Service {
	s := &Service{user_repo: user_repo}
	return s
}

// Operations
func (s *Service) GetAll() ([]User, error) {
	udbos, err := s.user_repo.GetAll()
	if err != nil {
		return nil, err
	}

	us := []User{}
	for _, udbo := range udbos {
		u := s.ToUser(&udbo)
		us = append(us, *u)
	}

	return us, nil
}
func (s *Service) GetByID(id uuid.UUID) (*User, error) {
	udbo, err := s.user_repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.ToUser(udbo), nil

}
func (s *Service) Create(user *User) (uuid.UUID, error) {
	// TODO: make sure user has a valid uuid before insert
	udbo := user.ToDBO()
	return s.user_repo.Create(udbo)
}

func (s *Service) ToUser(udb *userDBO) *User {
	return udb
}

func (u *User) ToDBO() *userDBO {
	return u
}
