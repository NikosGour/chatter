package user

import (
	"fmt"

	"github.com/NikosGour/chatter/internal/storage"
	"github.com/google/uuid"
)

type Repository interface {
	GetAll() ([]userDBO, error)
	GetByID(id uuid.UUID) (*userDBO, error)
	Create(user *userDBO) (uuid.UUID, error)
}

type repository struct {
	db *storage.PostgreSQLStorage
}

func NewRepository(db *storage.PostgreSQLStorage) Repository {
	ur := &repository{db: db}
	return ur
}

type userDBO = User

// Retrieves all user records from the database.
//
// Might return any sql error.
func (ur *repository) GetAll() ([]User, error) {
	udbos := []userDBO{}
	q := `SELECT id, username, password, date_created
		  FROM users`

	err := ur.db.Select(&udbos, q)
	if err != nil {
		return nil, err
	}
	return udbos, nil
}

// Retrieves a user given the UUID.
//
// Might return ErrGroupNotFound or any other sql error
func (ur *repository) GetByID(id uuid.UUID) (*User, error) {
	udbo := userDBO{}
	q := `SELECT id, username, password, date_created
		  FROM users
	      WHERE id = $1`

	err := ur.db.Get(&udbo, q, id)
	if err != nil {
		return nil, fmt.Errorf("on q=`%s`: %w", q, err)
	}

	return &udbo, nil
}

// Inserts a userinto a database.
//
// Returns the UUID of the created user.
// Might return any sql error
func (ur *repository) Create(user *userDBO) (uuid.UUID, error) {
	q := `INSERT INTO users (id, username, password, date_created)
		  VALUES (:id, :username, :password, :date_created)
		  RETURNING id;`

	insert_id := uuid.Nil
	stmt, err := ur.db.PrepareNamed(q)
	if err != nil {
		return uuid.Nil, fmt.Errorf("On q=`%s`: %w", q, err)
	}
	defer stmt.Close()

	err = stmt.Get(&insert_id, user)
	if err != nil {
		return uuid.Nil, fmt.Errorf("On q=`%s`: %w", q, err)
	}

	return insert_id, nil
}
