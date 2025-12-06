package channel

import (
	"fmt"

	"github.com/NikosGour/chatter/internal/storage"
	"github.com/google/uuid"
)

type Repository interface {
	GetAll() ([]channelDBO, error)
	GetByID(id uuid.UUID) (*channelDBO, error)
	Create(id uuid.UUID, chtype ChannelType) (uuid.UUID, error)
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

func (chr *repository) Create(id uuid.UUID, chtype ChannelType) (uuid.UUID, error) {
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
