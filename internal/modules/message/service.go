package message

import (
	"github.com/google/uuid"
)

type Service struct {
	message_repo Repository
}

func NewService(message_repo Repository) *Service {
	s := &Service{message_repo: message_repo}
	return s
}

// Retrieves all message records from the database.
//
// Might return any sql error.
func (s *Service) GetAll() ([]Message, error) {
	message_dbos, err := s.message_repo.GetAll()
	if err != nil {
		return nil, err
	}

	messages := []Message{}
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
func (s *Service) GetByID(id int64) (*Message, error) {
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

// Inserts a message into a database.
//
// Returns the id of the created message.
// Might return any sql error
func (s *Service) Create(message *Message) (uuid.UUID, error) {
	message_dbo := message.toDBO()
	return s.message_repo.Create(message_dbo)
}

// Transforms a message DBO to a message model
func (s *Service) toMessage(message_dbo messageDBO) (*Message, error) {
	message := &Message{
		Id:       message_dbo.Id,
		DateSent: message_dbo.DateSent,
	}

	panic("unimplemented")
	return message, nil
}

func (m *Message) toDBO() *messageDBO {
	mdbo := &messageDBO{
		Id:          m.Id,
		SenderId:    m.Sender.GetId(),
		RecipientId: m.Recipient.GetId(),
		DateSent:    m.DateSent,
	}
	return mdbo
}
