package group

import (
	"fmt"

	"github.com/NikosGour/chatter/internal/modules/channel"
	"github.com/NikosGour/chatter/internal/storage"
	"github.com/google/uuid"
)

type Repository interface {
	GetAll() ([]Group, error)
	GetByID(id uuid.UUID) (*Group, error)
	Create(group *Group) (uuid.UUID, error)
}

type repository struct {
	db *storage.PostgreSQLStorage

	chr channel.Repository
}

func NewRepository(db *storage.PostgreSQLStorage, chr channel.Repository) Repository {
	gr := &repository{db: db, chr: chr}
	return gr
}

type groupDBO = Group

func (gr *repository) GetAll() ([]Group, error) {
	gdbos := []groupDBO{}
	q := `SELECT id, name, date_created
	      FROM groups`

	err := gr.db.Select(&gdbos, q)
	if err != nil {
		return nil, err
	}

	gs := []Group{}
	for _, gdbo := range gdbos {
		g := gr.toGroup(&gdbo)
		gs = append(gs, *g)
	}

	return gs, nil
}

func (gr *repository) GetByID(id uuid.UUID) (*Group, error) {
	gdbo := groupDBO{}
	q := `SELECT id, name, date_created
		  FROM groups
	      WHERE id = $1`

	err := gr.db.Get(&gdbo, q, id)
	if err != nil {
		return nil, fmt.Errorf("on q=`%s`: %w", q, err)
	}

	g := gr.toGroup(&gdbo)
	return g, nil
}

func (gr *repository) Create(group *Group) (uuid.UUID, error) {
	id, err := gr.chr.Create(channel.ChannelTypeGroup)
	if err != nil {
		return uuid.Nil, fmt.Errorf("On channel create: %w", err)
	}

	group.Id = id

	gdbo := group.toDBO()
	q := `INSERT INTO groups (id, name, date_created)
		  VALUES (:id, :name, :date_created)
		  RETURNING id;`

	insert_id := uuid.Nil
	stmt, err := gr.db.PrepareNamed(q)
	if err != nil {
		return uuid.Nil, fmt.Errorf("On q=`%s`: %w", q, err)
	}
	defer stmt.Close()

	err = stmt.Get(&insert_id, gdbo)
	if err != nil {
		return uuid.Nil, fmt.Errorf("On q=`%s`: %w", q, err)
	}

	return insert_id, nil
}

func (gr *repository) toGroup(udb *groupDBO) *Group {
	return udb
}

func (g *Group) toDBO() *groupDBO {
	return g
}
