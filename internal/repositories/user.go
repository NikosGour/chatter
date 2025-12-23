package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/NikosGour/chatter/internal/models"
	"github.com/NikosGour/chatter/internal/storage"
	"github.com/google/uuid"
)

type UserRepository interface {
	GetAll() ([]UserDBO, error)
	GetByID(id uuid.UUID) (*UserDBO, error)
	GetByUsername(username string) ([]UserDBO, error)
	GetByTestUsername(username string) ([]UserDBO, error)
	Create(user *UserDBO) (uuid.UUID, error)
}

type userRepository struct {
	db *storage.PostgreSQLStorage
}

func NewUserRepository(db *storage.PostgreSQLStorage) UserRepository {
	ur := &userRepository{db: db}
	return ur
}

type UserDBO = models.User

// Retrieves all user records from the database.
//
// Might return any sql error.
func (ur *userRepository) GetAll() ([]UserDBO, error) {
	udbos := []UserDBO{}
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
func (ur *userRepository) GetByID(id uuid.UUID) (*UserDBO, error) {
	udbo := UserDBO{}
	q := `SELECT id, username, password, date_created
		  FROM users
	      WHERE id = $1`

	err := ur.db.Get(&udbo, q, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w:%s", models.ErrUserNotFound, id)
		}
		return nil, fmt.Errorf("on q=`%s`: %w", q, err)
	}

	return &udbo, nil
}
func (ur *userRepository) GetByUsername(username string) ([]UserDBO, error) {
	udbos := []UserDBO{}
	q := `SELECT id, username, password, date_created
		  FROM users
	      WHERE username = $1;`

	err := ur.db.Select(&udbos, q, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w:%s", models.ErrUserNotFound, username)
		}
		return nil, fmt.Errorf("on q=`%s`: %w", q, err)
	}

	return udbos, nil
}

func (ur *userRepository) GetByTestUsername(username string) ([]UserDBO, error) {
	udbos := []UserDBO{}
	q := `SELECT id, username, password, date_created
		  FROM users
	      WHERE username = $1 and is_test = true;`

	err := ur.db.Select(&udbos, q, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w:%s", models.ErrUserNotFound, username)
		}
		return nil, fmt.Errorf("on q=`%s`: %w", q, err)
	}

	return udbos, nil
}

// Inserts a userinto a database.
//
// Returns the UUID of the created user.
// Might return any sql error
func (ur *userRepository) Create(user *UserDBO) (uuid.UUID, error) {
	q := `INSERT INTO users (id, username, password, date_created, is_test)
		  VALUES (:id, :username, :password, :date_created, :is_test)
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
