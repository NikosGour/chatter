package channel

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/NikosGour/chatter/internal/storage"
	"github.com/google/uuid"
)

type Channel interface {
	GetId() uuid.UUID
}

type ChannelType string

const (
	ChannelTypeUser  ChannelType = "user"
	ChannelTypeGroup ChannelType = "group"
)

type Repository interface {
	GetAll() ([]channelDBO, error)
	GetByID(id uuid.UUID) (*channelDBO, error)
	Create(chtype ChannelType) (uuid.UUID, error)
}

type repository struct {
	db *storage.PostgreSQLStorage
}

func NewRepository(db *storage.PostgreSQLStorage) Repository {
	chr := &repository{db: db}
	return chr
}

type channelDBO struct {
	Id          uuid.UUID   `db:"id"`
	ChannelType ChannelType `db:"channel_type"`
}

func (chr *repository) GetAll() ([]channelDBO, error) {
	chdbos := []channelDBO{}
	q := `SELECT id, channel_type 
		  FROM channels`

	err := chr.db.Select(&chdbos, q)
	if err != nil {
		return nil, err
	}

	return chdbos, nil
}

func (chr *repository) GetByID(id uuid.UUID) (*channelDBO, error) {
	chdbo := channelDBO{}
	q := `SELECT id, channel_type 
		  FROM channels 
	      WHERE id = $1`

	err := chr.db.Get(&chdbo, q, id)
	if err != nil {
		return nil, fmt.Errorf("on q=`%s`: %w", q, err)
	}

	return &chdbo, nil
}

func (chr *repository) Create(chtype ChannelType) (uuid.UUID, error) {
	id, err := chr.createNewUUID()
	if err != nil {
		return uuid.Nil, fmt.Errorf("On createNewUUID: %w", err)
	}
	chdbo := channelDBO{
		Id:          id,
		ChannelType: chtype,
	}

	q := `INSERT INTO channels (id, channel_type)
		  VALUES (:id, :channel_type)
		  RETURNING id;`

	insert_id := uuid.Nil
	stmt, err := chr.db.PrepareNamed(q)
	if err != nil {
		return uuid.Nil, fmt.Errorf("On q=`%s`: %w", q, err)
	}
	defer stmt.Close()

	err = stmt.Get(&insert_id, chdbo)
	if err != nil {
		return uuid.Nil, fmt.Errorf("On q=`%s`: %w", q, err)
	}

	return insert_id, nil
}

func (chr *repository) createNewUUID() (uuid.UUID, error) {
	for {
		id := uuid.New()
		u, err := chr.GetByID(id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return id, nil
			}
			return uuid.Nil, fmt.Errorf("On GetById: %w", err)
		}
		if u == nil {
			return id, nil
		}
	}
}
