package services

import (
	"errors"
	"fmt"

	"github.com/NikosGour/chatter/internal/models"
	"github.com/NikosGour/chatter/internal/repositories"
	"github.com/google/uuid"
)

type TabService struct {
	tab_repo repositories.TabRepository
}

func NewTabService(tab_repo repositories.TabRepository) *TabService {
	s := &TabService{tab_repo: tab_repo}
	return s
}

// Operations
func (s *TabService) GetAll() ([]models.Tab, error) {
	tab_dbos, err := s.tab_repo.GetAll()
	if err != nil {
		return nil, err
	}

	tabs := []models.Tab{}
	for _, tab_dbo := range tab_dbos {
		tab := s.ToTab(&tab_dbo)
		tabs = append(tabs, *tab)
	}

	return tabs, nil
}
func (s *TabService) GetByID(id uuid.UUID) (*models.Tab, error) {
	tab_dbo, err := s.tab_repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.ToTab(tab_dbo), nil

}
func (s *TabService) GetByName(name string) ([]models.Tab, error) {
	tab_dbos, err := s.tab_repo.GetByName(name)
	if err != nil {
		return nil, err
	}

	tabs := []models.Tab{}
	for _, tab_dbo := range tab_dbos {
		tab := s.ToTab(&tab_dbo)
		tabs = append(tabs, *tab)
	}

	return tabs, nil
}

func (s *TabService) GetByServerID(server_id uuid.UUID) ([]models.Tab, error) {
	tab_dbos, err := s.tab_repo.GetByServerID(server_id)
	if err != nil {
		return nil, err
	}

	tabs := []models.Tab{}
	for _, tab_dbo := range tab_dbos {
		tab := s.ToTab(&tab_dbo)
		tabs = append(tabs, *tab)
	}

	return tabs, nil
}

func (s *TabService) Create(tab *models.Tab) (uuid.UUID, error) {
	id, err := s.generateUUID()
	if err != nil {
		return uuid.Nil, err
	}
	tab.Id = id

	tab_dbo := TabToDBO(tab)
	return s.tab_repo.Create(tab_dbo)
}

func (s *TabService) ToTab(tab_dbo *repositories.TabDBO) *models.Tab {
	return tab_dbo
}
func TabToDBO(tab *models.Tab) *repositories.TabDBO {
	return tab
}

func (s *TabService) generateUUID() (uuid.UUID, error) {
	id := uuid.New()

	for {
		ch, err := s.tab_repo.GetByID(id)
		if err != nil {
			if errors.Is(err, models.ErrTabNotFound) {
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
