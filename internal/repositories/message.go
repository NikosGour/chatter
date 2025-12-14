package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/NikosGour/chatter/internal/models"
	"github.com/NikosGour/chatter/internal/storage"
	"github.com/NikosGour/logging/log"
	"github.com/google/uuid"
)

type MessageRepository interface {
	GetAll() ([]MessageDBO, error)
	GetByID(id int64) (*MessageDBO, error)
	Create(group *MessageDBO) (int64, error)
}

type messageRepository struct {
	db *storage.PostgreSQLStorage
}

func NewMessageRepository(db *storage.PostgreSQLStorage) MessageRepository {
	mr := &messageRepository{db: db}
	return mr
}

type MessageDBO struct {
	Id          int64     `db:"id"`
	Text        string    `db:"text"`
	SenderId    uuid.UUID `db:"sender_id"`
	RecipientId uuid.UUID `db:"recipient_id"`
	DateSent    time.Time `db:"date_sent"`
}

// Retrieves all message records from the database.
//
// Might return any sql error.
func (mr *messageRepository) GetAll() ([]MessageDBO, error) {
	mdbos := []MessageDBO{}
	q := `SELECT id, "text", sender_id, recipient_id, date_sent
		  FROM messages`

	err := mr.db.Select(&mdbos, q)
	if err != nil {
		return nil, err
	}

	return mdbos, nil
}

// Retrieves a message given the id.
//
// Might return ErrGroupNotFound or any other sql error
func (mr *messageRepository) GetByID(id int64) (*MessageDBO, error) {
	mdbo := MessageDBO{}
	q := `SELECT id, "text", sender_id, recipient_id, date_sent
		  FROM messages
	      WHERE id = $1`

	err := mr.db.Get(&mdbo, q, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s:%d", models.ErrMessageNotFound, id)
		}

		msg := fmt.Errorf("on q=`%s`,id=`%d`: %w", q, id, err)
		log.Error("%s", msg)
		return nil, msg
	}

	return &mdbo, err
}

// Inserts a message into a database.
//
// Returns the id of the created message.
// Might return any sql error
func (mr *messageRepository) Create(message_dbo *MessageDBO) (int64, error) {
	q := `INSERT INTO messages ("text", sender_id, recipient_id, date_sent)
		  VALUES (:text, :sender_id, :recipient_id, :date_sent)
		  RETURNING id;`

	insert_id := int64(0)
	stmt, err := mr.db.PrepareNamed(q)
	if err != nil {
		return 0, fmt.Errorf("On q=`%s`: %w", q, err)
	}
	defer stmt.Close()

	err = stmt.Get(&insert_id, message_dbo)
	if err != nil {
		return 0, fmt.Errorf("On q=`%s`: %w", q, err)
	}

	return insert_id, nil
}
