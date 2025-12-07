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

type GroupRepository interface {
	GetAll() ([]GroupDBO, error)
	GetByID(id uuid.UUID) (*GroupDBO, error)
	Create(group *GroupDBO) (uuid.UUID, error)
	AddUserToGroup(user_id uuid.UUID, group_id uuid.UUID) error
	GetUsers(group_id uuid.UUID) ([]uuid.UUID, error)
}

type groupRepository struct {
	db *storage.PostgreSQLStorage
}

func NewGroupRepository(db *storage.PostgreSQLStorage) GroupRepository {
	gr := &groupRepository{db: db}
	return gr
}

type GroupDBO = models.Group

// Retrieves all group records from the database.
//
// Might return any sql error.
func (gr *groupRepository) GetAll() ([]GroupDBO, error) {
	gdbos := []GroupDBO{}
	q := `SELECT id, name, date_created
	      FROM groups`

	err := gr.db.Select(&gdbos, q)
	if err != nil {
		return nil, err
	}

	return gdbos, nil
}

// Retrieves a group given the UUID.
//
// Might return ErrGroupNotFound or any other sql error
func (gr *groupRepository) GetByID(id uuid.UUID) (*GroupDBO, error) {
	gdbo := GroupDBO{}
	q := `SELECT id, name, date_created
		  FROM groups
	      WHERE id = $1`

	err := gr.db.Get(&gdbo, q, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrGroupNotFound
		}

		msg := fmt.Errorf("on q=`%s`,id=`%s`: %w", q, id, err)
		log.Error("%s", msg)
		return nil, msg
	}

	return &gdbo, err
}

// Inserts a group into a database.
//
// Returns the UUID of the created group.
// Might return any sql error
func (gr *groupRepository) Create(group *GroupDBO) (uuid.UUID, error) {
	q := `INSERT INTO groups (id, name, date_created)
		  VALUES (:id, :name, :date_created)
		  RETURNING id;`

	insert_id := uuid.Nil
	stmt, err := gr.db.PrepareNamed(q)
	if err != nil {
		return uuid.Nil, fmt.Errorf("On q=`%s`: %w", q, err)
	}
	defer stmt.Close()

	err = stmt.Get(&insert_id, group)
	if err != nil {
		return uuid.Nil, fmt.Errorf("On q=`%s`: %w", q, err)
	}

	return insert_id, nil
}

// Adds a the user of the given UUID to the list of subscribed users of the group
//
// Might return ErrGroupNotFound or any other sql error
func (gr *groupRepository) AddUserToGroup(user_id uuid.UUID, group_id uuid.UUID) error {
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

// Get all the user UUIDs from a group's user list
//
// Might return ErrGroupHasNoUsers or any other sql error
func (gr *groupRepository) GetUsers(group_id uuid.UUID) ([]uuid.UUID, error) {
	user_ids := []uuid.UUID{}
	q := `SELECT user_id
		  FROM group_members
		  where group_id = $1;`
	err := gr.db.Select(&user_ids, q, group_id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrGroupHasNoUsers
		}
		return nil, err
	}

	return user_ids, nil
}
