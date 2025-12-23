package services

import (
	"github.com/NikosGour/chatter/internal/models"
	"github.com/NikosGour/chatter/internal/repositories"
	"github.com/google/uuid"
)

type MessageDTO = models.Message
type MessageService struct {
	message_repo repositories.MessageRepository

	tab_service *TabService
}

func NewMessageService(message_repo repositories.MessageRepository, tab_service *TabService) *MessageService {
	s := &MessageService{message_repo: message_repo, tab_service: tab_service}
	return s
}

// Retrieves all message records from the database.
//
// Might return any sql error.
func (s *MessageService) GetAll() ([]models.Message, error) {
	message_dbos, err := s.message_repo.GetAll()
	if err != nil {
		return nil, err
	}

	messages := []models.Message{}
	for _, message_dbo := range message_dbos {
		message, err := s.toMessage(message_dbo)
		if err != nil {
			return nil, err
		}
		messages = append(messages, *message)
	}
	return messages, nil
}

// Retrieves a message given the id.
//
// Might return ErrGroupNotFound or any other sql error
func (s *MessageService) GetByID(id int64) (*models.Message, error) {
	message_dbo, err := s.message_repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	message, err := s.toMessage(*message_dbo)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (s *MessageService) GetByTabID(tab_id uuid.UUID) ([]models.Message, error) {
	_, err := s.tab_service.GetByID(tab_id)
	if err != nil {
		return nil, err
	}

	message_dbos, err := s.message_repo.GetByTabID(tab_id)
	if err != nil {
		return nil, err
	}

	messages := []models.Message{}
	for _, message_dbo := range message_dbos {
		message, err := s.toMessage(message_dbo)
		if err != nil {
			return nil, err
		}
		messages = append(messages, *message)
	}
	return messages, nil
}

// Inserts a message into a database.
//
// Returns the id of the created message.
// Might return any sql error
func (s *MessageService) Create(message *models.Message) (int64, error) {
	message_dbo := messageToDBO(message)
	return s.message_repo.Create(message_dbo)
}

// Transforms a message DBO to a message model
func (s *MessageService) toMessage(message_dbo repositories.MessageDBO) (*models.Message, error) {
	message := &models.Message{
		Id:       message_dbo.Id,
		Text:     message_dbo.Text,
		DateSent: message_dbo.DateSent,
	}
	message.Sender = message_dbo.User
	message.Tab = message_dbo.Tab
	return message, nil
}

func (s *MessageService) MessageToDTO(m *models.Message) *MessageDTO {

	return m
}

func messageToDBO(m *models.Message) *repositories.MessageDBO {
	mdbo := &repositories.MessageDBO{
		Id:       m.Id,
		Text:     m.Text,
		SenderId: m.Sender.Id,
		TabId:    m.Tab.Id,
		DateSent: m.DateSent,
	}
	return mdbo
}
