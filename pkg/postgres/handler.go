package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/teris-io/shortid"
)

type PostgresHandler struct {
	dsn string
	db  *sql.DB
}

func NewPostgresHandler(dsn string) *PostgresHandler {
	return &PostgresHandler{dsn: dsn}
}

func (pgh *PostgresHandler) IsConnected() bool {
	return pgh.db != nil
}

func (pgh *PostgresHandler) EstablishConnection() error {
	conn, err := sql.Open("postgres", pgh.dsn)
	if err != nil {
		return err
	}
	pgh.db = conn
	return nil
}

type PGRole struct {
	RoleName string
	Password string
}

func (pgh *PostgresHandler) CreateRandomRole() (*PGRole, error) {
	if !pgh.IsConnected() {
		return nil, fmt.Errorf("Handler is currently not connected to a database")
	}
	roleName, err := shortid.Generate()
	passwd := generatePassword(32)
	if err != nil {
		return nil, fmt.Errorf("Unable to generate id %+v", err)
	}
	stmt := fmt.Sprintf("CREATE ROLE \"%s\" LOGIN PASSWORD '%s'", roleName, passwd)
	fmt.Println(stmt)
	_, err = pgh.db.Exec(stmt)
	if err != nil {
		return nil, fmt.Errorf("Unable to create role %s - %+v", roleName, err)
	}
	return &PGRole{RoleName: roleName, Password: passwd}, nil
}

type PGDatabase struct {
	Name string
	Role PGRole
}

func (pgh *PostgresHandler) CreateDatabase(dbName string, role *PGRole) (*PGDatabase, error) {

	if !pgh.IsConnected() {
		return nil, fmt.Errorf("Handler is currently not connected to a database")
	}
	tx, err := pgh.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("Unable to create transaction while creating database %s for role %s - %+v", dbName, role, err)
	}
	stmt := fmt.Sprintf("CREATE DATABASE \"%s\" OWNER \"%s\"", dbName, role.RoleName)
	_, err = tx.Exec(stmt)
	if err != nil {
		defer tx.Rollback()
		return nil, fmt.Errorf("Unable to create database %s for role %s - %+v", dbName, role.RoleName, err)
	}
	stmt = fmt.Sprintf("REVOKE ALL PRIVILEGES ON DATABASE \"%s\" FROM PUBLIC; GRANT ALL PRIVILEGES ON DATABASE \"%s\" TO \"%s\";", dbName, dbName, role.RoleName)
	_, err = tx.Exec(stmt)
	if err != nil {
		defer tx.Rollback()
		return nil, fmt.Errorf("Unable to set permissions for database %s for role %s - %+v", dbName, role.RoleName, err)
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return &PGDatabase{Name: dbName, Role: *role}, nil
}
