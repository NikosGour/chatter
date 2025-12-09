package repositories

import (
	"fmt"

	"github.com/NikosGour/chatter/internal/models"
	"github.com/NikosGour/chatter/internal/storage"
	"github.com/google/uuid"
)

type ChannelDBO struct {
	Id          uuid.UUID          `validate:"required" db:"id"`
	ChannelType models.ChannelType `db:"channel_type"`
}

type ChannelRepository interface {
	GetAll() ([]ChannelDBO, error)
	GetByID(id uuid.UUID) (*ChannelDBO, error)
	Create(chdbo *ChannelDBO) (uuid.UUID, error)
}

type channelRepository struct {
	db *storage.PostgreSQLStorage
}

func NewChannelRepository(db *storage.PostgreSQLStorage) ChannelRepository {
	chr := &channelRepository{db: db}
	return chr
}
func (chr *channelRepository) GetAll() ([]ChannelDBO, error) {
	chdbos := []ChannelDBO{}
	q := `SELECT id, channel_type
		  FROM channels`

	err := chr.db.Select(&chdbos, q)
	if err != nil {
		return nil, err
	}
	return chdbos, nil
}
func (chr *channelRepository) GetByID(id uuid.UUID) (*ChannelDBO, error) {
	chdbo := ChannelDBO{}
	q := `SELECT id, channel_type
		  FROM channels 
	      WHERE id = $1`

	err := chr.db.Get(&chdbo, q, id)
	if err != nil {
		return nil, fmt.Errorf("on q=`%s`: %w", q, err)
	}

	return &chdbo, nil
}

func (chr *channelRepository) Create(chdbo *ChannelDBO) (uuid.UUID, error) {
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
