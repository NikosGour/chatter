package user

import (
	"fmt"

	"github.com/NikosGour/chatter/internal/modules/channel"
	"github.com/NikosGour/chatter/internal/storage"
	"github.com/google/uuid"
)

type Repository interface {
	GetAll() ([]User, error)
	GetByID(id uuid.UUID) (*User, error)
	Create(user *User) (uuid.UUID, error)
}

type repository struct {
	db *storage.PostgreSQLStorage

	chr channel.Repository
}

func NewRepository(db *storage.PostgreSQLStorage, chr channel.Repository) Repository {
	ur := &repository{db: db, chr: chr}
	return ur
}

type userDBO = User

func (ur *repository) GetAll() ([]User, error) {
	udbos := []userDBO{}
	q := `SELECT id, username, password, date_created
		  FROM users`

	err := ur.db.Select(&udbos, q)
	if err != nil {
		return nil, err
	}

	us := []User{}
	for _, udbo := range udbos {
		u := ur.toUser(&udbo)
		us = append(us, *u)
	}

	return us, nil
}

func (ur *repository) GetByID(id uuid.UUID) (*User, error) {
	udbo := userDBO{}
	q := `SELECT id, username, password, date_created
		  FROM users
	      WHERE id = $1`

	err := ur.db.Get(&udbo, q, id)
	if err != nil {
		return nil, fmt.Errorf("on q=`%s`: %w", q, err)
	}

	u := ur.toUser(&udbo)
	return u, nil
}

func (ur *repository) Create(user *User) (uuid.UUID, error) {
	id, err := ur.chr.Create(channel.ChannelTypeUser)
	if err != nil {
		return uuid.Nil, fmt.Errorf("On channel create: %w", err)
	}

	user.Id = id

	udbo := user.toDBO()
	q := `INSERT INTO users (id, username, password, date_created)
		  VALUES (:id, :username, :password, :date_created)
		  RETURNING id;`

	insert_id := uuid.Nil
	stmt, err := ur.db.PrepareNamed(q)
	if err != nil {
		return uuid.Nil, fmt.Errorf("On q=`%s`: %w", q, err)
	}
	defer stmt.Close()

	err = stmt.Get(&insert_id, udbo)
	if err != nil {
		return uuid.Nil, fmt.Errorf("On q=`%s`: %w", q, err)
	}

	return insert_id, nil
}

func (ur *repository) toUser(udb *userDBO) *User {
	return udb
}

func (u *User) toDBO() *userDBO {
	return u
}
