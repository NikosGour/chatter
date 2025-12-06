package message

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/NikosGour/chatter/internal/storage"
	"github.com/NikosGour/logging/log"
	"github.com/google/uuid"
)

type Repository interface {
	GetAll() ([]messageDBO, error)
	GetByID(id int64) (*messageDBO, error)
	Create(group *messageDBO) (uuid.UUID, error)
}

type repository struct {
	db *storage.PostgreSQLStorage
}

func NewRepository(db *storage.PostgreSQLStorage) Repository {
	mr := &repository{db: db}
	return mr
}

type messageDBO struct {
	Id          int64     `db:"id"`
	SenderId    uuid.UUID `db:"sender_id"`
	RecipientId uuid.UUID `db:"recipient_id"`
	DateSent    time.Time `db:"date_sent"`
}

// Retrieves all message records from the database.
//
// Might return any sql error.
func (mr *repository) GetAll() ([]messageDBO, error) {
	mdbos := []messageDBO{}
	q := `SELECT id, sender_id, recipient_id, date_sent
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
func (mr *repository) GetByID(id int64) (*messageDBO, error) {
	mdbo := messageDBO{}
	q := `SELECT id, sender_id, recipient_id, date_sent
		  FROM messages
	      WHERE id = $1`

	err := mr.db.Get(&mdbo, q, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMessageNotFound
		}

		msg := fmt.Errorf("on q=`%s`,id=`%s`: %w", q, id, err)
		log.Error("%s", msg)
		return nil, msg
	}

	return &mdbo, err
}

// Inserts a message into a database.
//
// Returns the id of the created message.
// Might return any sql error
func (mr *repository) Create(message_dbo *messageDBO) (uuid.UUID, error) {
	q := `INSERT INTO messages (name, date_created)
		  VALUES (:name, :date_created)
		  RETURNING id;`

	insert_id := uuid.Nil
	stmt, err := mr.db.PrepareNamed(q)
	if err != nil {
		return uuid.Nil, fmt.Errorf("On q=`%s`: %w", q, err)
	}
	defer stmt.Close()

	err = stmt.Get(&insert_id, message_dbo)
	if err != nil {
		return uuid.Nil, fmt.Errorf("On q=`%s`: %w", q, err)
	}

	return insert_id, nil
}
