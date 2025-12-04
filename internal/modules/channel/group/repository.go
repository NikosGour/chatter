package group

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/NikosGour/chatter/internal/modules/channel"
	"github.com/NikosGour/chatter/internal/modules/channel/user"
	"github.com/NikosGour/chatter/internal/storage"
	"github.com/google/uuid"
)

type Repository interface {
	GetAll() ([]Group, error)
	GetByID(id uuid.UUID) (*Group, error)
	Create(group *Group) (uuid.UUID, error)
	AddUserToGroup(user_id uuid.UUID, group_id uuid.UUID) error
}

type repository struct {
	db *storage.PostgreSQLStorage

	channel_repo channel.Repository
	user_repo    user.Repository
}

func NewRepository(db *storage.PostgreSQLStorage, channel_repo channel.Repository, user_repo user.Repository) Repository {
	gr := &repository{db: db, channel_repo: channel_repo, user_repo: user_repo}
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
	id, err := gr.channel_repo.Create(channel.ChannelTypeGroup)
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

func (gr *repository) AddUserToGroup(user_id uuid.UUID, group_id uuid.UUID) error {
	_, err := gr.user_repo.GetByID(user_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user.ErrUserNotFound
		}
		return err
	}

	_, err = gr.GetByID(group_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrGroupNotFound
		}
		return err
	}

	q := `INSERT INTO group_members (group_id, user_id)
		  VALUES (:group,:user)`

	_, err = gr.db.NamedExec(q, struct {
		Group uuid.UUID `db:"group"`
		User  uuid.UUID `db:"user"`
	}{Group: group_id, User: user_id})
	if err != nil {
		return fmt.Errorf("On q=`%s`: %w", q, err)
	}

	return nil
}

func (gr *repository) toGroup(udb *groupDBO) *Group {
	return udb
}

func (g *Group) toDBO() *groupDBO {
	return g
}
