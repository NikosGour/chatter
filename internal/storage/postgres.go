package storage

import (
	"fmt"

	"github.com/NikosGour/chatter/internal/common"
	"github.com/NikosGour/chatter/internal/projectpath"
	"github.com/NikosGour/logging/log"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgreSQLStorage struct {
	*sqlx.DB
}

func NewPostgreSQLStorage() *PostgreSQLStorage {
	p := &PostgreSQLStorage{}
	p.init_database()
	return p
}

func (st *PostgreSQLStorage) init_database() {
	var (
		host   = common.Dotenv[common.EnvPOSTGRES_HOST_ADDRESS]
		port   = common.Dotenv[common.EnvPOSTGRES_PORT]
		user   = common.Dotenv[common.EnvPOSTGRES_USER]
		dbpass = common.Dotenv[common.EnvPOSTGRES_ROOT_PASSWORD]
		dbname = common.Dotenv[common.EnvPOSTGRES_DB]
	)

	conn_string := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable", host, port, user, dbpass, dbname)

	var err error
	st.DB, err = sqlx.Connect("postgres", conn_string)
	if err != nil {
		log.Fatal("%s", err)
	}

	err = st.CreateTables()
	if err != nil {
		log.Fatal("%s", err)
	}
}

func (st *PostgreSQLStorage) CreateTables() error {
	_, err := sqlx.LoadFile(st, projectpath.RootFile("db/create_channels.sql"))
	if err != nil {
		return fmt.Errorf("on LoadFile(create_channels): %w", err)
	}
	_, err = sqlx.LoadFile(st, projectpath.RootFile("db/create_users.sql"))
	if err != nil {
		return fmt.Errorf("on LoadFile(create_users): %w", err)
	}
	_, err = sqlx.LoadFile(st, projectpath.RootFile("db/create_groups.sql"))
	if err != nil {
		return fmt.Errorf("on LoadFile(create_groups): %w", err)
	}
	_, err = sqlx.LoadFile(st, projectpath.RootFile("db/create_group_members.sql"))
	if err != nil {
		return fmt.Errorf("on LoadFile(create_group_members): %w", err)
	}
	_, err = sqlx.LoadFile(st, projectpath.RootFile("db/create_messages.sql"))
	if err != nil {
		return fmt.Errorf("on LoadFile(create_messages): %w", err)
	}

	return nil
}

func (st *PostgreSQLStorage) DropTables() error {
	_, err := sqlx.LoadFile(st, projectpath.RootFile("db/drop_all.sql"))
	if err != nil {
		return fmt.Errorf("on LoadFile(drop_all): %w", err)
	}

	return nil

}
