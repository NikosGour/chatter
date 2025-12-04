package group

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/NikosGour/chatter/internal/modules/channel"
	"github.com/NikosGour/chatter/internal/modules/channel/user"
	"github.com/NikosGour/chatter/internal/storage"
	"github.com/NikosGour/logging/log"
	"github.com/google/uuid"
)

type Repository interface {
	GetAll() ([]groupDBO, error)
	GetByID(id uuid.UUID) (*groupDBO, error)
	Create(group *Group) (uuid.UUID, error)
	AddUserToGroup(user_id uuid.UUID, group_id uuid.UUID) error
	GetUsers(group_id uuid.UUID) ([]uuid.UUID, error)
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

func (gr *repository) GetAll() ([]groupDBO, error) {
	gdbos := []groupDBO{}
	q := `SELECT id, name, date_created
	      FROM groups`

	err := gr.db.Select(&gdbos, q)
	if err != nil {
		return nil, err
	}

	return gdbos, nil
}

func (gr *repository) GetByID(id uuid.UUID) (*groupDBO, error) {
	gdbo := groupDBO{}
	q := `SELECT id, name, date_created
		  FROM groups
	      WHERE id = $1`

	err := gr.db.Get(&gdbo, q, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrGroupNotFound
		}

		msg := fmt.Errorf("on q=`%s`,id=`%s`: %w", q, id, err)
		log.Error("%s", msg)
		return nil, msg
	}

	return &gdbo, err
}

func (gr *repository) Create(group *Group) (uuid.UUID, error) {
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
	q := `INSERT INTO group_members (group_id, user_id)
		  VALUES (:group,:user)`

	_, err := gr.db.NamedExec(q, struct {
		Group uuid.UUID `db:"group"`
		User  uuid.UUID `db:"user"`
	}{Group: group_id, User: user_id})
	if err != nil {
		return fmt.Errorf("On q=`%s`: %w", q, err)
	}

	return nil
}

func (gr *repository) GetUsers(group_id uuid.UUID) ([]uuid.UUID, error) {
	user_ids := []uuid.UUID{}
	q := `SELECT user_id
		  FROM group_members
		  where group_id = $1;`
	err := gr.db.Select(&user_ids, q, group_id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrGroupHasNoUsers
		}
		return nil, err
	}

	return user_ids, nil
}
