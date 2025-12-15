package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/NikosGour/chatter/internal/models"
	"github.com/NikosGour/chatter/internal/storage"
	"github.com/google/uuid"
)

type TabRepository interface {
	GetAll() ([]TabDBO, error)
	GetByID(id uuid.UUID) (*TabDBO, error)
	GetByServerID(server_id uuid.UUID) ([]TabDBO, error)
	Create(Tab *TabDBO) (uuid.UUID, error)
}

type tabRepository struct {
	db *storage.PostgreSQLStorage
}

func NewTabRepository(db *storage.PostgreSQLStorage) TabRepository {
	tr := &tabRepository{db: db}
	return tr
}

type TabDBO = models.Tab

// Retrieves all Tab records from the database.
//
// Might return any sql error.
func (tr *tabRepository) GetAll() ([]TabDBO, error) {
	tab_dbos := []TabDBO{}
	// q := `SELECT id, name, server_id, date_created
	// 	  FROM tabs;`

	q := `SELECT t.*,s.id as "server.id", s.name as "server.name"
		  FROM tabs t
          JOIN servers s ON t.server_id = s.id`

	err := tr.db.Select(&tab_dbos, q)
	if err != nil {
		return nil, fmt.Errorf("on q=`%s`: %w", q, err)
	}
	return tab_dbos, nil
}

// Retrieves a tab given the UUID.
//
// Might return ErrTabNotFound or any other sql error
func (tr *tabRepository) GetByID(id uuid.UUID) (*TabDBO, error) {
	tab_dbo := TabDBO{}
	q := `SELECT id, name, server_id, date_created
		  FROM tabs
	      WHERE id = $1`

	err := tr.db.Get(&tab_dbo, q, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w:%s", models.ErrTabNotFound, id)
		}
		return nil, fmt.Errorf("on q=`%s`: %w", q, err)
	}

	return &tab_dbo, nil
}

// Retrieves a tab given the server UUID.
//
// Might return ErrTabNotFound or any other sql error
func (tr *tabRepository) GetByServerID(server_id uuid.UUID) ([]TabDBO, error) {
	tab_dbos := []TabDBO{}
	q := `SELECT id, name, server_id, date_created
		  FROM tabs
	      WHERE server_id = $1`

	err := tr.db.Select(&tab_dbos, q, server_id)
	if err != nil {
		return nil, fmt.Errorf("on q=`%s`: %w", q, err)
	}

	return tab_dbos, nil
}

// Inserts a tab into a database.
//
// Returns the UUID of the created Tab.
// Might return any sql error
func (tr *tabRepository) Create(Tab *TabDBO) (uuid.UUID, error) {
	q := `INSERT INTO Tabs (id, name, server_id, date_created)
		  VALUES (:id, :name, :server_id, :date_created)
		  RETURNING id;`

	insert_id := uuid.Nil
	stmt, err := tr.db.PrepareNamed(q)
	if err != nil {
		return uuid.Nil, fmt.Errorf("On q=`%s`: %w", q, err)
	}
	defer stmt.Close()

	err = stmt.Get(&insert_id, Tab)
	if err != nil {
		return uuid.Nil, fmt.Errorf("On q=`%s`: %w", q, err)
	}

	return insert_id, nil
}
