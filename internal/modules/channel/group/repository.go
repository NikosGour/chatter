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
	GetAll() ([]Group, error)
	GetByID(id uuid.UUID) (*Group, error)
	Create(group *Group) (uuid.UUID, error)
	AddUserToGroup(user_id uuid.UUID, group_id uuid.UUID) error
	GetUsers(group_id uuid.UUID) ([]user.User, error)
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
		g, err := gr.toGroup(&gdbo)
		if err != nil {
			return nil, err
		}
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrGroupNotFound
		}

		msg := fmt.Errorf("on q=`%s`,id=`%s`: %w", q, id, err)
		log.Error("%s", msg)
		return nil, msg
	}

	g, err := gr.toGroup(&gdbo)
	if err != nil {
		return nil, err
	}
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

func (gr *repository) GetUsers(group_id uuid.UUID) ([]user.User, error) {
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

	us := []user.User{}
	for _, user_id := range user_ids {
		udbo, err := gr.user_repo.GetByID(user_id)
		if err != nil {
			if errors.Is(err, user.ErrUserNotFound) {
				log.Warn("while getting users for group: %s, tried to get missing user: %s", group_id, user_id)
			}
			return nil, err
		}
		u := gr.user_repo.ToUser(udbo)
		us = append(us, *u)
	}

	return us, nil
}

func (gr *repository) toGroup(udb *groupDBO) (*Group, error) {
	users, err := gr.GetUsers(udb.Id)
	if err != nil {
		if errors.Is(err, ErrGroupHasNoUsers) {
			users = []user.User{}
		} else {
			return nil, err
		}
	}
	udb.Users = users
	return udb, nil
}

func (g *Group) toDBO() *groupDBO {
	return g
}
