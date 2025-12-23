package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/NikosGour/chatter/internal/models"
	"github.com/NikosGour/chatter/internal/storage"
	"github.com/NikosGour/logging/log"
	"github.com/google/uuid"
)

type ServerRepository interface {
	GetAll() ([]ServerDBO, error)
	GetByID(id uuid.UUID) (*ServerDBO, error)
	GetByName(name string) ([]ServerDBO, error)
	GetByTestName(name string) ([]ServerDBO, error)
	Create(Server *ServerDBO) (uuid.UUID, error)
	AddUserToServer(user_id uuid.UUID, server_id uuid.UUID) error
	GetUsers(server_id uuid.UUID) ([]uuid.UUID, error)
}

type serverRepository struct {
	db *storage.PostgreSQLStorage
}

func NewServerRepository(db *storage.PostgreSQLStorage) ServerRepository {
	sr := &serverRepository{db: db}
	return sr
}

type ServerDBO = models.Server

// Retrieves all servers from the database.
//
// Might return any sql error.
func (sr *serverRepository) GetAll() ([]ServerDBO, error) {
	server_dbos := []ServerDBO{}
	q := `SELECT id, name, date_created
	      FROM servers;`

	err := sr.db.Select(&server_dbos, q)
	if err != nil {
		return nil, err
	}

	return server_dbos, nil
}

// Retrieves a server given the UUID.
//
// Might return ErrServerNotFound or any other sql error
func (sr *serverRepository) GetByID(id uuid.UUID) (*ServerDBO, error) {
	server_dbo := ServerDBO{}
	q := `SELECT id, name, date_created
		  FROM servers
	      WHERE id = $1;`

	err := sr.db.Get(&server_dbo, q, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w:%s", models.ErrServerNotFound, id)
		}

		msg := fmt.Errorf("on q=`%s`,id=`%s`: %w", q, id, err)
		log.Error("%s", msg)
		return nil, msg
	}

	return &server_dbo, err
}

func (sr *serverRepository) GetByName(name string) ([]ServerDBO, error) {
	sdbos := []ServerDBO{}
	q := `SELECT id, name, date_created
		  FROM servers
	      WHERE name = $1;`

	err := sr.db.Select(&sdbos, q, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w:%s", models.ErrUserNotFound, name)
		}
		return nil, fmt.Errorf("on q=`%s`: %w", q, err)
	}

	return sdbos, nil
}
func (sr *serverRepository) GetByTestName(name string) ([]ServerDBO, error) {
	sdbos := []ServerDBO{}
	q := `SELECT id, name, date_created
		  FROM servers
	      WHERE name = $1 and is_test = true;`

	err := sr.db.Select(&sdbos, q, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w:%s", models.ErrUserNotFound, name)
		}
		return nil, fmt.Errorf("on q=`%s`: %w", q, err)
	}

	return sdbos, nil

}

// Inserts a server into a database.
//
// Returns the UUID of the created server.
// Might return any sql error
func (sr *serverRepository) Create(server *ServerDBO) (uuid.UUID, error) {
	q := `INSERT INTO servers (id, name, date_created, is_test)
		  VALUES (:id, :name, :date_created, :is_test)
		  RETURNING id;`

	insert_id := uuid.Nil
	stmt, err := sr.db.PrepareNamed(q)
	if err != nil {
		return uuid.Nil, fmt.Errorf("On q=`%s`: %w", q, err)
	}
	defer stmt.Close()

	err = stmt.Get(&insert_id, server)
	if err != nil {
		return uuid.Nil, fmt.Errorf("On q=`%s`: %w", q, err)
	}

	return insert_id, nil
}

// Adds a the user of the given UUID to the list of subscribed users of the server
//
// Might return ErrServerNotFound or any other sql error
func (sr *serverRepository) AddUserToServer(user_id uuid.UUID, server_id uuid.UUID) error {
	q := `INSERT INTO server_members (server_id, user_id)
		  VALUES (:server,:user)`

	_, err := sr.db.NamedExec(q, struct {
		Server uuid.UUID `db:"server"`
		User   uuid.UUID `db:"user"`
	}{Server: server_id, User: user_id})
	if err != nil {
		return fmt.Errorf("On q=`%s`: %w", q, err)
	}

	return nil
}

// Get all the user UUIDs from a server's user list
//
// Might return ErrServerHasNoUsers or any other sql error
func (sr *serverRepository) GetUsers(server_id uuid.UUID) ([]uuid.UUID, error) {
	user_ids := []uuid.UUID{}
	q := `SELECT user_id
		  FROM server_members
		  where server_id = $1;`
	err := sr.db.Select(&user_ids, q, server_id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w:%s", models.ErrServerHasNoUsers, server_id)
		}
		return nil, err
	}

	return user_ids, nil
}
